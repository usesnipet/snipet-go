import { mount, unmount } from 'svelte';
import App from './App.svelte';
import { WidgetState, type WidgetPosition } from './widget-state.svelte';

const HOST_ID = 'snipet-widget-root';

export type SnipetWidgetApi = {
	open: () => void;
	close: () => void;
	toggle: () => void;
	destroy: () => void;
};

type EmbedConfig = {
	clientCode: string;
	position: WidgetPosition;
};

type Runtime = {
	host: HTMLElement;
	app: ReturnType<typeof mount>;
	state: WidgetState;
};

let runtime: Runtime | null = null;

function resolveScript(): HTMLScriptElement | null {
	if (document.currentScript instanceof HTMLScriptElement) {
		return document.currentScript;
	}

	return (
		document.querySelector<HTMLScriptElement>('script[data-snipet-widget]') ??
		document.querySelector<HTMLScriptElement>('script[src*="snipet-widget"]') ??
		document.querySelector<HTMLScriptElement>('script[data-client][src*="embed"]')
	);
}

function readConfig(script: HTMLScriptElement | null): EmbedConfig {
	const clientCode = script?.dataset.client?.trim() ?? '';
	if (!clientCode) {
		throw new Error('[SnipetWidget] data-client is required on the script tag');
	}

	const rawPosition = script?.dataset.position?.trim() ?? 'bottom-right';
	const position: WidgetPosition =
		rawPosition === 'bottom-left' ? 'bottom-left' : 'bottom-right';

	return { clientCode, position };
}

function createApi(state: WidgetState): SnipetWidgetApi {
	return {
		open: state.openPanel,
		close: state.closePanel,
		toggle: state.toggle,
		destroy,
	};
}

function destroy(): void {
	if (!runtime) return;

	void unmount(runtime.app);
	runtime.host.remove();
	runtime = null;
}

function mountWidget(config: EmbedConfig): SnipetWidgetApi {
	if (runtime) {
		return createApi(runtime.state);
	}

	const existing = document.getElementById(HOST_ID);
	if (existing) {
		existing.remove();
	}

	const host = document.createElement('div');
	host.id = HOST_ID;
	host.style.cssText = 'all:initial;position:fixed;z-index:2147483647;';
	document.body.appendChild(host);

	const shadow = host.attachShadow({ mode: 'open' });
	const state = new WidgetState();
	const app = mount(App, {
		target: shadow,
		props: {
			clientCode: config.clientCode,
			position: config.position,
			state,
		},
	});

	runtime = { host, app, state };
	return createApi(state);
}

function init(): SnipetWidgetApi {
	const config = readConfig(resolveScript());

	if (!document.body) {
		let resolved: SnipetWidgetApi | null = null;
		const deferred: SnipetWidgetApi = {
			open: () => resolved?.open(),
			close: () => resolved?.close(),
			toggle: () => resolved?.toggle(),
			destroy: () => resolved?.destroy(),
		};

		document.addEventListener(
			'DOMContentLoaded',
			() => {
				resolved = mountWidget(config);
				window.SnipetWidget = resolved;
			},
			{ once: true },
		);

		return deferred;
	}

	return mountWidget(config);
}

declare global {
	interface Window {
		SnipetWidget: SnipetWidgetApi;
	}
}

const api = init();
window.SnipetWidget = api;

export default api;

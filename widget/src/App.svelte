<script lang="ts">
	import type { WidgetPosition, WidgetState } from './widget-state.svelte';

	interface Props {
		clientCode: string;
		position?: WidgetPosition;
		state: WidgetState;
	}

	let { clientCode, position = 'bottom-right', state }: Props = $props();
</script>

<div class={['widget', position === 'bottom-left' && 'left']} data-client={clientCode}>
	<div
		class={['panel', state.open && 'open']}
		role="dialog"
		aria-label="Snipet chat"
		aria-hidden={!state.open}
	>
		<header class="header">
			<span class="brand">Snipet</span>
			<button type="button" class="icon-btn" aria-label="Fechar" onclick={state.closePanel}>
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
					<path d="M18 6 6 18M6 6l12 12" />
				</svg>
			</button>
		</header>
		<div class="body">
			<p class="placeholder">Chat em breve</p>
		</div>
	</div>

	<button
		type="button"
		class="launcher"
		aria-label={state.open ? 'Fechar chat' : 'Abrir chat'}
		aria-expanded={state.open}
		onclick={state.toggle}
	>
		{#if state.open}
			<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
				<path d="M18 6 6 18M6 6l12 12" />
			</svg>
		{:else}
			<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
				<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
			</svg>
		{/if}
	</button>
</div>

<style>
	:host {
		all: initial;
		font-family: 'Segoe UI', system-ui, sans-serif;
		line-height: 1.4;
		color: #1a1a1a;
	}

	.widget {
		position: fixed;
		bottom: 20px;
		right: 20px;
		z-index: 2147483646;
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 12px;
		box-sizing: border-box;
	}

	.widget.left {
		right: auto;
		left: 20px;
		align-items: flex-start;
	}

	.panel {
		width: min(360px, calc(100vw - 40px));
		height: min(520px, calc(100vh - 100px));
		display: flex;
		flex-direction: column;
		background: #fff;
		border: 1px solid #e5e5e5;
		border-radius: 16px;
		box-shadow: 0 12px 40px rgba(0, 0, 0, 0.14);
		overflow: hidden;
		opacity: 0;
		transform: translateY(12px) scale(0.96);
		pointer-events: none;
		transition:
			opacity 180ms ease,
			transform 180ms ease;
	}

	.panel.open {
		opacity: 1;
		transform: translateY(0) scale(1);
		pointer-events: auto;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 14px 16px;
		border-bottom: 1px solid #eee;
		background: #fafafa;
	}

	.brand {
		font-size: 15px;
		font-weight: 650;
		letter-spacing: -0.02em;
	}

	.icon-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		padding: 0;
		border: none;
		border-radius: 8px;
		background: transparent;
		color: #555;
		cursor: pointer;
	}

	.icon-btn:hover {
		background: #eee;
		color: #111;
	}

	.body {
		flex: 1;
		display: grid;
		place-items: center;
		padding: 24px;
		background: #fff;
	}

	.placeholder {
		margin: 0;
		color: #888;
		font-size: 14px;
	}

	.launcher {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 56px;
		height: 56px;
		padding: 0;
		border: none;
		border-radius: 50%;
		background: #111;
		color: #fff;
		cursor: pointer;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
		transition:
			transform 160ms ease,
			background 160ms ease;
	}

	.launcher:hover {
		background: #2a2a2a;
		transform: scale(1.04);
	}

	.launcher:active {
		transform: scale(0.98);
	}
</style>

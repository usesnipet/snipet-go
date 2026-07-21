export type WidgetPosition = 'bottom-right' | 'bottom-left';

export class WidgetState {
	open = $state(false);

	openPanel = (): void => {
		this.open = true;
	};

	closePanel = (): void => {
		this.open = false;
	};

	toggle = (): void => {
		this.open = !this.open;
	};
}

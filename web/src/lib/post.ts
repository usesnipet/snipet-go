import { ON_INPUT, ON_CHANGE } from '@sjsf/form';
import type { Schema, UiSchemaRoot } from '@sjsf/form';
import type { FromSchema } from 'json-schema-to-ts';

export const schema = {
	title: 'Post',
	type: 'object',
	properties: {
		title: { title: 'Title', type: 'string' },
		content: { title: 'Content', type: 'string', minLength: 10 }
	},
	required: ['title', 'content']
} as const satisfies Schema;

export type Model = FromSchema<typeof schema>;

export const uiSchema = {
	content: { 'ui:components': { textWidget: 'textareaWidget' } }
} as const satisfies UiSchemaRoot;

export const initialValue = { title: 'New post', content: '' };
export const fieldsValidationMode = ON_INPUT | ON_CHANGE;

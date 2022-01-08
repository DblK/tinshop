import { dev } from '$app/env';

export const variables = {
    basePath: dev ? '/' : '/admin/'
};
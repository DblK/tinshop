import { dev } from '$app/env';

export const variables = {
    basePath: dev ? '/admin/' : '/admin/',
    api: dev ? 'http://localhost:3000/api' : '/api'
};
export class HTTPTransport {
    constructor(private url: string) {}

    async request(method: string, params: unknown): Promise<unknown> {
        const jsonrpc = {
            jsonrpc: '2.0',
            method,
            params,
            id: 1,
        }
        const body = JSON.stringify(jsonrpc)
        try {
            const result = await fetch(this.url, {
                method: 'POST',
                body,
                headers: {
                    'Content-Type': 'application/json',
                },
            })
            const json = await result.json()
            if (json.error) {
                throw json
            }
            return json.result
        } catch (e) {
            console.log('HTTPTransport encountered error:', e)
            throw e
        }
    }
}

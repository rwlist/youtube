export class RpcError {
    code: number
    message: string
    data: unknown

    constructor(obj: unknown) {
        if (
            typeof obj === 'object' &&
            obj != null &&
            'code' in obj &&
            'message' in obj &&
            'data' in obj
        ) {
            const e = obj as RpcError
            this.code = e.code
            this.message = e.message
            this.data = e.data
            Object.assign(this, e)
            return
        }

        throw new Error('invalid rpc error')
    }
}

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
                throw new RpcError(json.error)
            }
            return json.result
        } catch (e) {
            console.log('HTTPTransport encountered error:', e)
            throw e
        }
    }
}

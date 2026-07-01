export type ApiRequest = {
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  token?: string;
  route: string
};

type response = {
  ok: boolean,
  data: any,
  status: number,
  err: string
}

export async function reqServer(value: ApiRequest) {
  let headers = {
    "Content-Type": "application/json",
    ...(
      value.token && {
        "Authorization": `Bearer ${value.token}`
      }
    )
  }
  let payload: RequestInit = {
    method: value.method,
    headers: headers,
    ...(value.method != "GET" && {
      body: JSON.stringify(value.body)
    })
  }
  try {
    const res = await fetch(
      `${import.meta.env.VITE_SERVER_URL}/${value.route}`,
      payload
    )
    const data = res.json()
    const T: response = {
      ok: res.ok,
      data: data,
      status: res.status,
      err: ""
    }
    return T
  }
  catch (err) {
    const T: response = {
      ok: false,
      data: null,
      status: 0,
      err: err instanceof Error ? err.message : String(err),
    }
    return T
  }
}

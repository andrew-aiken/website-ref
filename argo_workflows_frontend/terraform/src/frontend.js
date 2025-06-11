function handler(event) {
    let request = event.request;

    const subdomain = request.headers.host.value.split('.')[0];

    if (!/\..+/.test(request.uri)) {
        request.uri = `/${subdomain}/index.html`;
    } else {
        request.uri = `/${subdomain}${request.uri}`;
    }
    return request;
}

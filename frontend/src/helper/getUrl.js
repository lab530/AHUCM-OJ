export function getCurrentURL() {
    // https://www.example.com/path/page?param1=value1&param2=value2#section1
    // return https://www.example.com/path/page?param1=value1&param2=value2#section1。
    const currentURL = window.location.href;
    return currentURL;
}

export function getCurrentQueryParams() {
    // https://www.example.com/path/page?param1=value1&param2=value2
    // return { param1: 'value1', param2: 'value2' }
    const searchParams = new URLSearchParams(window.location.search);
    const queryParams = Object.fromEntries(searchParams.entries());
    return queryParams;
}

export function getPathName() {
    // https://www.example.com/path/page?param1=value1&param2=value2
    // return /path/page?param1=value1&param2=value2
    // 获取完整的URL  
    const fullUrl = window.location.href;

    // 创建一个URL对象，这样我们可以更容易地解析URL的各个部分  
    const url = new URL(fullUrl, window.location.origin);

    // 提取pathname，即主机域名之后的部分  
    const pathAfterHostname = url.pathname + url.search;

    // 返回结果  
    return pathAfterHostname;
}

export function getProblemIdFromURL(param) {
    // 获取 param 后面的 value
    const urlParams = new URLSearchParams(window.location.search);
    const paramValue = urlParams.get(param);
    return paramValue;
}
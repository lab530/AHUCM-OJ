export function formatCurrentDateTime() {
    const currentDate = new Date();
    currentDate.setHours(currentDate.getHours() - 8);
    const year = currentDate.getFullYear();
    const month = String(currentDate.getMonth() + 1).padStart(2, '0');
    const day = String(currentDate.getDate()).padStart(2, '0');
    const hours = String(currentDate.getHours()).padStart(2, '0');
    const minutes = String(currentDate.getMinutes()).padStart(2, '0');
    const seconds = String(currentDate.getSeconds()).padStart(2, '0');
    const milliseconds = String(currentDate.getMilliseconds()).padStart(3, '0');

    const formattedDateTime = `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
    return formattedDateTime;
}
export function convertToBackendTime(timeString) {
    const date = new Date(timeString);
    const year = date.getUTCFullYear();
    const month = String(date.getUTCMonth() + 1).padStart(2, '0');
    const day = String(date.getUTCDate()).padStart(2, '0');
    const hours = String(date.getUTCHours()).padStart(2, '0');
    const minutes = String(date.getUTCMinutes()).padStart(2, '0');
    const seconds = String(date.getUTCSeconds()).padStart(2, '0');
    const milliseconds = String(date.getUTCMilliseconds()).padStart(3, '0');

    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`;
}

export function getDuration(start, end) {
    const startTime = new Date(start);
    const endTime = new Date(end);
    const duration = endTime.getTime() - startTime.getTime();

    const days = Math.floor(duration / (1000 * 60 * 60 * 24));
    const hours = Math.floor((duration % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((duration % (1000 * 60 * 60)) / (1000 * 60));

    return `${days}天${hours}时${minutes}分`;
}

export function getCurrentTime() {
    return new Date();
}

export function formatDate(date) {
    // 格式化日期的方法  
    const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' };
    return new Date(date).toLocaleString('zh-CN', options);
}

export function getTotalDuration(startDate, endDate) {
    // 计算总赛时的方法  
    const start = new Date(startDate);
    const end = new Date(endDate);
    const duration = end - start;
    const days = Math.floor(duration / (1000 * 60 * 60 * 24));
    const hours = Math.floor((duration % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((duration % (1000 * 60 * 60)) / (1000 * 60));
    return `${days}天${hours}时${minutes}分`;
}

export function isBeforeStart(item) {
    // 检查是否在开始前  
    const now = new Date();
    return now < new Date(item.start_at);
}

export function isDuring(item) {
    // 检查是否正在进行  
    const now = new Date();
    return new Date(item.start_at) <= now && now <= new Date(item.end_at);
}  
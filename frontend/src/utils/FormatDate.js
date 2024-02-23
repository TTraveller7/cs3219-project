export function formatDateTimestampToMonthDayTime(date) {
    const month = date.getMonth();
    const day = date.getDate();
    const hour24 = date.getHours();
    let hour12 = hour24 % 12;
    const minute = date.getMinutes();
    const second = date.getSeconds();
    let suffix;

    if (hour24 / 12 >= 1) {
        suffix = 'pm'
    } else {
        suffix = 'am'
    }

    return day + "/" + month + " " + hour12 + ":" + minute + ":" + second + " " + suffix;
}
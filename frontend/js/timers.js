export function getRemainingTime(expiresAtMs) {
    const now = new Date().getTime();
    const distance = expiresAtMs - now;
    return distance > 0 ? distance : 0;
}

export function formatTime(ms) {
    if (ms <= 0) return '0m 0s';
    const minutes = Math.floor(ms / (1000 * 60));
    const seconds = Math.floor((ms % (1000 * 60)) / 1000);
    return `${minutes}m ${seconds}s`;
}

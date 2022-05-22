export function countLines(str: string): number {
    return str.split(/\r\n|\r|\n/).length;
}

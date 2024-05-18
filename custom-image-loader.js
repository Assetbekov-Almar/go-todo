import path from 'path';
import { basePath } from './next.config';

export default function myImageLoader({ src }) {
    if (basePath && path.isAbsolute(src)) {
        return `${basePath}${src}`;
    }
    return src;
}

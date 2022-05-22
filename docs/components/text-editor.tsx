import { useRef, useState } from 'react';
import { countLines } from '../util';

const TextEditor: React.FC = () => {
    const [sourceCode, setSourceCode] = useState('');
    const lines = useRef<HTMLDivElement>(null);

    return (
        <div className="flex h-full w-full flex-row overflow-hidden font-mono">
            <div ref={lines} className="w-10 bg-yellow-50 pr-2 text-right text-slate-400">
                {[...Array(countLines(sourceCode))].map((x, i) => (
                    <div key={i}>{i + 1}</div>
                ))}
            </div>
            <textarea
                className="flex-grow resize-none bg-yellow-50 text-slate-900 outline-none"
                value={sourceCode}
                onChange={e => setSourceCode(e.target.value)}
                onScroll={e => {
                    if (lines.current) lines.current.style.marginTop = `-${(e.target as HTMLElement).scrollTop}px`;
                }}
                autoCapitalize="off"
                autoCorrect="off"
                autoComplete="off"
                spellCheck="false"
                wrap="off"
            ></textarea>
        </div>
    );
};

export default TextEditor;

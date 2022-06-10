import { Fragment, useCallback, useRef, useState } from 'react';
import { countLines } from '../util';
import { Menu } from '@headlessui/react';
import { FaChevronDown } from 'react-icons/fa';
import examples from './examples';

const endpoint = 'https://api.aspen.junlim.dev/run';

const CodeEditor: React.FC = () => {
    const [sourceCode, setSourceCode] = useState(examples[0].code);
    const [output, setOutput] = useState('');
    const lines = useRef<HTMLDivElement>(null);
    const textarea = useRef<HTMLTextAreaElement>(null);

    const onRun = useCallback(() => {
        setOutput('Waiting for remote server...');
        fetch(endpoint, {
            method: 'POST',
            body: sourceCode,
        })
            .then(response => response.text())
            .then(out => setOutput(out))
            .catch(_ => setOutput('Internal server error.'));
    }, [sourceCode]);

    return (
        <Fragment>
            <div className="mb-4 flex flex-row items-center">
                <h1 className="mr-auto text-2xl font-medium">The Aspen Playground</h1>
                <button
                    className="mr-4 rounded-md bg-blue-500 px-3 py-1 font-medium text-white hover:bg-blue-600"
                    onClick={onRun}
                >
                    Run
                </button>
                <Menu as="div" className="relative mr-2">
                    <Menu.Button className="flex flex-row items-center border border-neutral-300 bg-yellow-50 py-1 px-3 font-medium hover:bg-yellow-100">
                        Examples
                        <FaChevronDown className="ml-2" />
                    </Menu.Button>
                    <Menu.Items className="absolute right-0 mt-2 w-56 rounded bg-white py-3 shadow-[0_0_6px_2px_rgba(0,0,0,0.1)]">
                        {examples.map(example => (
                            <Menu.Item key={example.name}>
                                {({ active }) => (
                                    <button
                                        onClick={() => setSourceCode(example.code)}
                                        className={`${
                                            active && 'bg-blue-500 text-white'
                                        } block w-full px-3 py-1 text-left`}
                                    >
                                        {example.name}
                                    </button>
                                )}
                            </Menu.Item>
                        ))}
                    </Menu.Items>
                </Menu>
            </div>
            <div className="divide-y divide-neutral-200 border border-neutral-200">
                <div className="flex h-120 w-full flex-row overflow-hidden font-mono">
                    <div ref={lines} className="w-10 bg-yellow-50 pr-2 text-right text-slate-400">
                        {[...Array(countLines(sourceCode))].map((x, i) => (
                            <div key={i}>{i + 1}</div>
                        ))}
                    </div>
                    <textarea
                        ref={textarea}
                        className="flex-grow resize-none bg-yellow-50 text-slate-900 outline-none"
                        value={sourceCode}
                        onChange={e => setSourceCode(e.target.value)}
                        onScroll={e => {
                            if (lines.current)
                                lines.current.style.marginTop = `-${(e.target as HTMLElement).scrollTop}px`;
                        }}
                        onKeyDown={e => {
                            // see https://stackoverflow.com/questions/40331780/reactjs-handle-tab-character-in-textarea
                            if (e.key === 'Tab' && !e.shiftKey) {
                                e.preventDefault();
                                const value = textarea.current!.value;
                                const selectionStart = textarea.current!.selectionStart;
                                const selectionEnd = textarea.current!.selectionEnd;

                                const newText =
                                    value.substring(0, selectionStart) + '    ' + value.substring(selectionEnd);
                                textarea.current!.value = newText; // need to set this otherwise changing the selection wont work

                                textarea.current!.selectionStart = selectionEnd + 4 - (selectionEnd - selectionStart);
                                textarea.current!.selectionEnd = selectionEnd + 4 - (selectionEnd - selectionStart);

                                setSourceCode(newText);
                            }
                        }}
                        autoCapitalize="off"
                        autoCorrect="off"
                        autoComplete="off"
                        spellCheck="false"
                        wrap="off"
                    ></textarea>
                </div>
                <div className="min-h-[15rem] overflow-x-auto whitespace-pre bg-neutral-100 p-4 font-mono">
                    {output}
                </div>
            </div>
        </Fragment>
    );
};

export default CodeEditor;

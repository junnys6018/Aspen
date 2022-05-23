import type { NextPage } from 'next';
import Head from 'next/head';
import CodeEditor from '../components/code-editor';

const Home: NextPage = () => {
    return (
        <div>
            <Head>
                <title>Aspen</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <div className="relative mx-auto mt-12" style={{ width: '800px' }}>
                <CodeEditor />
            </div>
        </div>
    );
};

export default Home;

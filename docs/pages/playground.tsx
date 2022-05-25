import type { NextPage } from 'next';
import Head from 'next/head';
import CodeEditor from '../components/code-editor';
import Footer from '../components/footer';
import Header from '../components/header';

const Playground: NextPage = () => {
    return (
        <div>
            <Head>
                <title>Aspen</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <Header />
            <div className="container relative mx-auto my-8">
                <CodeEditor />
            </div>
            <Footer />
        </div>
    );
};

export default Playground;

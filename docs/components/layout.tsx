import React from 'react';
import Head from 'next/head';
import Footer from './footer';
import Header from './header';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return (
        <div>
            <Head>
                <title>Aspen</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <Header />
            <main className="container relative mx-auto my-8">{children}</main>
            <Footer />
        </div>
    );
};

export default Layout;

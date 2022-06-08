import React from 'react';
import Head from 'next/head';
import Footer from './footer';
import Header from './header';
import SideNav from './side-nav';
import table from '../table';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return (
        <div>
            <Head>
                <title>Aspen</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <Header />
            <main className="container relative mx-auto my-8">
                <div className="flex flex-col md:flex-row">
                    <SideNav table={table} />
                    <div className="flex-grow">{children}</div>
                </div>
            </main>
            <Footer />
        </div>
    );
};

export default Layout;

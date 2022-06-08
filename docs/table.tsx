export type Table = {
    header: string;
    items: {
        name: string;
        slug?: string;
        href?: string;
    }[];
}[];

const table: Table = [
    {
        header: 'Welcome!',
        items: [
            {
                name: 'Home',
                slug: '/',
            },
            {
                name: 'Playground',
                slug: '/playground',
            },
            {
                name: 'GitHub',
                href: 'https://github.com/junnys6018/Aspen',
            },
        ],
    },
    {
        header: 'Basics',
        items: [
            {
                name: 'Installation',
                slug: '/d',
            },
            {
                name: 'GitHub',
                slug: '/e',
            },
            {
                name: 'Patterns',
                slug: '/f',
            },
        ],
    },
];

export default table;

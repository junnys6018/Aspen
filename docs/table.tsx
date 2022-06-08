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
        header: 'Guides',
        items: [
            {
                name: 'Basics',
                slug: '/basics',
            },
            {
                name: 'Types',
                slug: '/types',
            },
            {
                name: 'Control Flow',
                slug: '/control-flow',
            },
            {
                name: 'Functions',
                slug: '/functions',
            },
        ],
    },
    {
        header: 'Reference',
        items: [
            {
                name: 'Aspen CLI',
                slug: '/cli',
            },
            {
                name: 'Grammar',
                slug: '/grammar',
            },
        ],
    },
];

export default table;

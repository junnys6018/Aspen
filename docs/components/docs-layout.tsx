const DocsLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return (
        <article className="prose max-w-none prose-headings:text-blue-500 prose-h1:font-semibold prose-h2:font-semibold">
            {children}
        </article>
    );
};

export default DocsLayout;

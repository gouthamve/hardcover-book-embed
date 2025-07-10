(function() {
    // Hardcover Book Widget
    const WIDGET_VERSION = '1.0.0';
    
    // Default configuration
    const defaultConfig = {
        apiUrl: 'http://localhost:8080',
        username: null,
        maxWidth: '800px',
        columns: 'auto-fill',
        minColumnWidth: '120px',
        gap: '1rem',
        showPoweredBy: true
    };

    // Widget styles
    const styles = `
        .hardcover-widget {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: var(--hw-max-width, 800px);
            margin: 0 auto;
            padding: var(--hw-padding, 1rem);
        }

        .hardcover-widget * {
            box-sizing: border-box;
        }

        .hw-books-grid {
            display: grid;
            grid-template-columns: repeat(var(--hw-columns, auto-fill), minmax(var(--hw-min-column-width, 120px), 1fr));
            gap: var(--hw-gap, 1rem);
            padding: 0;
            margin: 0;
            list-style: none;
        }

        .hw-book-item {
            position: relative;
            transition: all 0.3s ease;
        }

        .hw-book-item:hover {
            transform: translateY(-4px);
        }

        .hw-book-link {
            display: block;
            text-decoration: none;
            color: inherit;
        }

        .hw-book-cover {
            position: relative;
            width: 100%;
            padding-bottom: 150%; /* 2:3 aspect ratio */
            background: #e5e7eb;
            border-radius: 6px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
            transition: all 0.3s ease;
        }

        .hw-book-item:hover .hw-book-cover {
            box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
        }

        .hw-book-cover img {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
        }

        .hw-book-title-overlay {
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            background: linear-gradient(to top, rgba(0, 0, 0, 0.9) 0%, rgba(0, 0, 0, 0.7) 50%, transparent 100%);
            color: white;
            padding: 2rem 0.75rem 0.75rem;
            opacity: 0;
            transition: all 0.3s ease;
            font-size: 0.875rem;
            font-weight: 600;
            line-height: 1.3;
        }

        .hw-book-item:hover .hw-book-title-overlay {
            opacity: 1;
        }

        .hw-loading {
            text-align: center;
            padding: 3rem;
            color: #6b7280;
        }

        .hw-error {
            text-align: center;
            padding: 2rem;
            color: #dc2626;
            background: #fef2f2;
            border-radius: 8px;
            margin: 1rem;
        }

        .hw-empty-state {
            text-align: center;
            padding: 3rem;
            color: #6b7280;
        }

        .hw-powered-by {
            margin-top: 1.5rem;
            text-align: center;
            font-size: 0.75rem;
            color: #6b7280;
        }

        .hw-powered-by a {
            color: #2563eb;
            text-decoration: none;
            transition: all 0.3s ease;
        }

        .hw-powered-by a:hover {
            text-decoration: underline;
        }

        @media (max-width: 600px) {
            .hw-books-grid {
                grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
                gap: 0.75rem;
            }
            
            .hw-book-title-overlay {
                font-size: 0.75rem;
                padding: 1.5rem 0.5rem 0.5rem;
            }
        }

        @media (max-width: 400px) {
            .hw-books-grid {
                grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
                gap: 0.5rem;
            }
        }
    `;

    // Main widget class
    class HardcoverWidget {
        constructor(element, config = {}) {
            this.element = element;
            this.config = { ...defaultConfig, ...config };
            this.init();
        }

        init() {
            // Add widget class
            this.element.classList.add('hardcover-widget');
            
            // Set CSS variables
            this.element.style.setProperty('--hw-max-width', this.config.maxWidth);
            this.element.style.setProperty('--hw-columns', this.config.columns);
            this.element.style.setProperty('--hw-min-column-width', this.config.minColumnWidth);
            this.element.style.setProperty('--hw-gap', this.config.gap);
            
            // Load books
            this.loadBooks();
        }

        async loadBooks() {
            if (!this.config.username) {
                this.showError('Username is required. Please add data-username attribute.');
                return;
            }
            
            this.showLoading();
            
            try {
                const response = await fetch(`${this.config.apiUrl}/api/books/currently-reading/${this.config.username}`);
                
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`);
                }
                
                const data = await response.json();
                
                if (data.count === 0) {
                    this.showEmptyState();
                } else {
                    this.renderBooks(data.books);
                }
                
                // Store reference for programmatic access
                this.element._hardcoverWidget = this;
                
                // Dispatch success event
                this.element.dispatchEvent(new CustomEvent('hardcover:loaded', {
                    detail: { count: data.count, books: data.books },
                    bubbles: true
                }));
            } catch (error) {
                console.error('Hardcover Widget Error:', error);
                this.showError('Failed to load books. Please try again later.');
                
                // Dispatch error event
                this.element.dispatchEvent(new CustomEvent('hardcover:error', {
                    detail: { error: error.message },
                    bubbles: true
                }));
            }
        }

        showLoading() {
            this.element.innerHTML = '<div class="hw-loading">Loading books...</div>';
        }

        showError(message) {
            this.element.innerHTML = `<div class="hw-error">${message}</div>`;
        }

        showEmptyState() {
            this.element.innerHTML = '<div class="hw-empty-state">No books currently being read</div>';
        }

        renderBooks(books) {
            const booksHtml = books.map(book => this.renderBook(book)).join('');
            
            let html = `<ul class="hw-books-grid">${booksHtml}</ul>`;
            
            if (this.config.showPoweredBy) {
                html += `
                    <div class="hw-powered-by">
                        <a href="https://hardcover.app/@${this.config.username}" target="_blank" rel="noopener">Currently reading on Hardcover</a>
                    </div>
                `;
            }
            
            this.element.innerHTML = html;
        }

        renderBook(book) {
            const cover = book.book.image && book.book.image.url
                ? `<img src="${book.book.image.url}" alt="${book.book.title} cover" loading="lazy">`
                : '';
            
            const bookUrl = `https://hardcover.app/books/${book.book.slug}`;

            return `
                <li class="hw-book-item">
                    <a href="${bookUrl}" target="_blank" rel="noopener" class="hw-book-link">
                        <div class="hw-book-cover">
                            ${cover}
                            <div class="hw-book-title-overlay">${book.book.title}</div>
                        </div>
                    </a>
                </li>
            `;
        }
    }

    // Add styles to page
    function addStyles() {
        if (document.getElementById('hardcover-widget-styles')) return;
        
        const styleSheet = document.createElement('style');
        styleSheet.id = 'hardcover-widget-styles';
        styleSheet.textContent = styles;
        document.head.appendChild(styleSheet);
    }

    // Auto-initialize widgets
    function autoInit() {
        addStyles();
        
        // Find all elements with data-hardcover-widget attribute
        const widgets = document.querySelectorAll('[data-hardcover-widget]');
        
        widgets.forEach(element => {
            // Parse config from data attributes
            const config = {};
            
            if (element.dataset.apiUrl) config.apiUrl = element.dataset.apiUrl;
            if (element.dataset.username) config.username = element.dataset.username;
            if (element.dataset.maxWidth) config.maxWidth = element.dataset.maxWidth;
            if (element.dataset.columns) config.columns = element.dataset.columns;
            if (element.dataset.minColumnWidth) config.minColumnWidth = element.dataset.minColumnWidth;
            if (element.dataset.gap) config.gap = element.dataset.gap;
            if (element.dataset.showPoweredBy !== undefined) {
                config.showPoweredBy = element.dataset.showPoweredBy !== 'false';
            }
            
            new HardcoverWidget(element, config);
        });
    }

    // Export for manual initialization
    window.HardcoverWidget = HardcoverWidget;
    window.HardcoverWidget.init = autoInit;
    window.HardcoverWidget.addStyles = addStyles;

    // Auto-initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', autoInit);
    } else {
        autoInit();
    }
})();
(function() {
    // Hardcover Review Widget
    const WIDGET_VERSION = '1.0.0';
    
    // Default configuration
    const defaultConfig = {
        apiUrl: 'http://localhost:8080',
        username: null,
        maxWidth: '800px',
        showPoweredBy: true,
        maxReviewLength: 300
    };

    // Widget styles
    const styles = `
        .hardcover-review-widget {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: var(--hrw-max-width, 800px);
            margin: 0 auto;
            padding: var(--hrw-padding, 1rem);
        }

        .hardcover-review-widget * {
            box-sizing: border-box;
        }

        .hrw-reviews-list {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            list-style: none;
            padding: 0;
            margin: 0;
        }

        .hrw-review-item {
            display: flex;
            gap: 1rem;
            padding: 1rem;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            transition: all 0.3s ease;
            background: #ffffff;
        }

        .hrw-review-item:hover {
            box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
            border-color: #d1d5db;
        }

        .hrw-book-cover {
            flex-shrink: 0;
            width: 80px;
            height: 120px;
            border-radius: 4px;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        .hrw-book-cover img {
            width: 100%;
            height: 100%;
            object-fit: cover;
        }

        .hrw-mobile-cover {
            display: none !important;
        }

        .hrw-review-content {
            flex: 1;
            min-width: 0;
        }

        .hrw-review-header {
            margin-bottom: 0.5rem;
        }

        .hrw-book-title {
            font-size: 1.125rem;
            font-weight: 600;
            color: #374151;
            margin: 0 0 0.25rem 0;
            text-decoration: none;
            display: inline-block;
        }

        .hrw-book-title:hover {
            color: #2563eb;
            text-decoration: underline;
        }

        .hrw-review-meta {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            font-size: 0.875rem;
            color: #6b7280;
        }

        .hrw-rating {
            display: flex;
            align-items: center;
            gap: 0.25rem;
        }

        .hrw-stars {
            display: inline-flex;
            gap: 2px;
        }

        .hrw-star {
            width: 16px;
            height: 16px;
            fill: #fbbf24;
        }

        .hrw-star.empty {
            fill: #e5e7eb;
        }

        .hrw-review-date {
            color: #6b7280;
        }

        .hrw-review-text {
            margin: 0.75rem 0 0 0;
            line-height: 1.6;
            color: #374151;
            font-size: 0.9375rem;
        }

        .hrw-spoiler-warning {
            display: inline-block;
            background: #fef3c7;
            color: #92400e;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.75rem;
            font-weight: 500;
            margin-bottom: 0.5rem;
        }

        .hrw-read-more {
            color: #2563eb;
            text-decoration: none;
            font-size: 0.875rem;
            font-weight: 500;
            display: inline-block;
            margin-top: 0.5rem;
        }

        .hrw-read-more:hover {
            text-decoration: underline;
        }

        .hrw-loading {
            text-align: center;
            padding: 3rem;
            color: #6b7280;
        }

        .hrw-error {
            text-align: center;
            padding: 2rem;
            color: #dc2626;
            background: #fef2f2;
            border-radius: 8px;
            margin: 1rem;
        }

        .hrw-empty-state {
            text-align: center;
            padding: 3rem;
            color: #6b7280;
        }

        .hrw-powered-by {
            margin-top: 1.5rem;
            text-align: center;
            font-size: 0.75rem;
            color: #6b7280;
        }

        .hrw-powered-by a {
            color: #2563eb;
            text-decoration: none;
            transition: all 0.3s ease;
        }

        .hrw-powered-by a:hover {
            text-decoration: underline;
        }

        @media (max-width: 600px) {
            .hrw-review-item {
                padding: 0.75rem;
            }

            .hrw-book-cover {
                width: 60px;
                height: 90px;
            }

            .hrw-review-header {
                display: flex;
                gap: 0.75rem;
                align-items: flex-start;
                margin-bottom: 0.75rem;
            }

            .hrw-review-header .hrw-book-cover {
                flex-shrink: 0;
            }

            .hrw-book-info {
                flex: 1;
                min-width: 0;
            }

            .hrw-book-title {
                font-size: 1rem;
                line-height: 1.3;
            }

            .hrw-review-meta {
                flex-wrap: wrap;
                font-size: 0.8125rem;
                margin-top: 0.25rem;
            }

            .hrw-review-text {
                font-size: 0.875rem;
            }

            .hrw-desktop-cover {
                display: none;
            }

            .hrw-mobile-cover {
                display: block !important;
            }
        }
    `;

    // Main widget class
    class HardcoverReviewWidget {
        constructor(element, config = {}) {
            this.element = element;
            this.config = { ...defaultConfig, ...config };
            this.init();
        }

        init() {
            // Add widget class
            this.element.classList.add('hardcover-review-widget');
            
            // Set CSS variables
            this.element.style.setProperty('--hrw-max-width', this.config.maxWidth);
            
            // Load reviews
            this.loadReviews();
        }

        async loadReviews() {
            if (!this.config.username) {
                this.showError('Username is required. Please add data-username attribute.');
                return;
            }
            
            this.showLoading();
            
            try {
                const endpoint = `/api/books/reviews/${this.config.username}`;
                const response = await fetch(`${this.config.apiUrl}${endpoint}`);
                
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}`);
                }
                
                const data = await response.json();
                
                if (data.count === 0) {
                    this.showEmptyState();
                } else {
                    this.renderReviews(data.books);
                }
                
                // Store reference for programmatic access
                this.element._hardcoverReviewWidget = this;
                
                // Dispatch success event
                this.element.dispatchEvent(new CustomEvent('hardcover:reviews-loaded', {
                    detail: { count: data.count, reviews: data.books },
                    bubbles: true
                }));
            } catch (error) {
                console.error('Hardcover Review Widget Error:', error);
                this.showError('Failed to load reviews. Please try again later.');
                
                // Dispatch error event
                this.element.dispatchEvent(new CustomEvent('hardcover:reviews-error', {
                    detail: { error: error.message },
                    bubbles: true
                }));
            }
        }

        showLoading() {
            this.element.innerHTML = '<div class="hrw-loading">Loading reviews...</div>';
        }

        showError(message) {
            this.element.innerHTML = `<div class="hrw-error">${message}</div>`;
        }

        showEmptyState() {
            this.element.innerHTML = '<div class="hrw-empty-state">No reviews yet</div>';
        }

        renderStars(rating) {
            if (!rating) return '';
            
            const fullStars = Math.floor(rating);
            const hasHalfStar = rating % 1 !== 0;
            const emptyStars = 5 - Math.ceil(rating);
            
            let starsHtml = '<span class="hrw-stars">';
            
            // Full stars
            for (let i = 0; i < fullStars; i++) {
                starsHtml += '<svg class="hrw-star" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>';
            }
            
            // Half star
            if (hasHalfStar) {
                starsHtml += '<svg class="hrw-star" viewBox="0 0 20 20"><defs><linearGradient id="halfstar-review"><stop offset="50%" stop-color="#fbbf24"/><stop offset="50%" stop-color="#e5e7eb"/></linearGradient></defs><path fill="url(#halfstar-review)" d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>';
            }
            
            // Empty stars
            for (let i = 0; i < emptyStars; i++) {
                starsHtml += '<svg class="hrw-star empty" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>';
            }
            
            starsHtml += '</span>';
            return starsHtml;
        }

        formatDate(dateObj) {
            if (!dateObj) return '';
            // Handle both string dates and object dates (from Go's Date type)
            const dateString = typeof dateObj === 'string' ? dateObj : dateObj.Time;
            if (!dateString) return '';
            const date = new Date(dateString);
            const options = { year: 'numeric', month: 'short', day: 'numeric' };
            return date.toLocaleDateString('en-US', options);
        }

        truncateText(text, maxLength) {
            if (!text || text.length <= maxLength) return text;
            return text.substring(0, maxLength).trim() + '...';
        }

        renderReviews(reviews) {
            const reviewsHtml = reviews.map(review => this.renderReview(review)).join('');
            
            let html = `<ul class="hrw-reviews-list">${reviewsHtml}</ul>`;
            
            if (this.config.showPoweredBy) {
                html += `
                    <div class="hrw-powered-by">
                        <a href="https://hardcover.app/@${this.config.username}" target="_blank" rel="noopener">Book reviews on Hardcover</a>
                    </div>
                `;
            }
            
            this.element.innerHTML = html;
        }

        renderReview(review) {
            const cover = review.book.image && review.book.image.url
                ? `<img src="${review.book.image.url}" alt="${review.book.title} cover" loading="lazy">`
                : '';
            
            const reviewUrl = `https://hardcover.app/books/${review.book.slug}/reviews/@${this.config.username}`;
            const reviewText = review.review_raw || '';
            const truncatedText = this.truncateText(reviewText, this.config.maxReviewLength);
            const needsReadMore = reviewText.length > this.config.maxReviewLength;

            return `
                <li class="hrw-review-item">
                    <div class="hrw-book-cover hrw-desktop-cover">
                        ${cover}
                    </div>
                    <div class="hrw-review-content">
                        <div class="hrw-review-header">
                            <div class="hrw-book-cover hrw-mobile-cover">
                                ${cover}
                            </div>
                            <div class="hrw-book-info">
                                <a href="${reviewUrl}" target="_blank" rel="noopener" class="hrw-book-title">${review.book.title}</a>
                                <div class="hrw-review-meta">
                                    ${review.rating ? `
                                        <div class="hrw-rating">
                                            ${this.renderStars(review.rating)}
                                        </div>
                                    ` : ''}
                                    ${review.reviewed_at ? `
                                        <span class="hrw-review-date">${this.formatDate(review.reviewed_at)}</span>
                                    ` : ''}
                                </div>
                            </div>
                        </div>
                        ${review.review_has_spoilers ? '<span class="hrw-spoiler-warning">Contains spoilers</span>' : ''}
                        ${reviewText ? `
                            <p class="hrw-review-text">${truncatedText}</p>
                            ${needsReadMore ? `<a href="${reviewUrl}" target="_blank" rel="noopener" class="hrw-read-more">Read full review →</a>` : ''}
                        ` : ''}
                    </div>
                </li>
            `;
        }
    }

    // Add styles to page
    function addStyles() {
        if (document.getElementById('hardcover-review-widget-styles')) return;
        
        const styleSheet = document.createElement('style');
        styleSheet.id = 'hardcover-review-widget-styles';
        styleSheet.textContent = styles;
        document.head.appendChild(styleSheet);
    }

    // Auto-initialize widgets
    function autoInit() {
        addStyles();
        
        // Find all elements with data-hardcover-review-widget attribute
        const widgets = document.querySelectorAll('[data-hardcover-review-widget]');
        
        widgets.forEach(element => {
            // Parse config from data attributes
            const config = {};
            
            if (element.dataset.apiUrl) config.apiUrl = element.dataset.apiUrl;
            if (element.dataset.username) config.username = element.dataset.username;
            if (element.dataset.maxWidth) config.maxWidth = element.dataset.maxWidth;
            if (element.dataset.maxReviewLength) config.maxReviewLength = parseInt(element.dataset.maxReviewLength);
            if (element.dataset.showPoweredBy !== undefined) {
                config.showPoweredBy = element.dataset.showPoweredBy !== 'false';
            }
            
            new HardcoverReviewWidget(element, config);
        });
    }

    // Export for manual initialization
    window.HardcoverReviewWidget = HardcoverReviewWidget;
    window.HardcoverReviewWidget.init = autoInit;
    window.HardcoverReviewWidget.addStyles = addStyles;

    // Auto-initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', autoInit);
    } else {
        autoInit();
    }
})();
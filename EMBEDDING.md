# Hardcover Book Widget - Embedding Guide

This guide explains how to embed any Hardcover user's currently reading or last read books on any website.

## Quick Start

The simplest way to embed your books is to add a single line of HTML to your website:

```html
<!-- Add this where you want the books to appear -->
<div data-hardcover-widget data-api-url="http://localhost:8080" data-username="your-username"></div>

<!-- Add this before closing </body> tag -->
<script src="http://localhost:8080/widget.js"></script>
```

That's it! Replace "your-username" with any Hardcover username to show their currently reading books.

To show last read books instead, add `data-book-type="last-read"`:

```html
<div data-hardcover-widget data-api-url="http://localhost:8080" data-username="your-username" data-book-type="last-read"></div>
<script src="http://localhost:8080/widget.js"></script>
```

## Embedding Methods

### Method 1: Simple HTML (Recommended)

Add these two lines to your HTML:

```html
<div data-hardcover-widget data-api-url="http://localhost:8080" data-username="your-username"></div>
<script src="http://localhost:8080/widget.js"></script>
```

### Method 2: Iframe

For complete isolation from your site's styles:

```html
<iframe 
    src="http://localhost:8080/embed.html?username=your-username" 
    width="100%" 
    height="400"
    frameborder="0"
    style="border: none;">
</iframe>

<!-- For last read books -->
<iframe 
    src="http://localhost:8080/embed.html?username=your-username&type=last-read" 
    width="100%" 
    height="400"
    frameborder="0"
    style="border: none;">
</iframe>
```

### Method 3: JavaScript API

For more control over initialization:

```html
<div id="my-books"></div>
<script src="http://localhost:8080/widget.js"></script>
<script>
    // Manual initialization
    const widget = new HardcoverWidget(
        document.getElementById('my-books'),
        {
            apiUrl: 'http://localhost:8080',
            username: 'your-username',
            bookType: 'currently-reading', // or 'last-read'
            maxWidth: '600px',
            showPoweredBy: true
        }
    );
</script>
```

## Configuration Options

You can customize the widget using data attributes:

```html
<div 
    data-hardcover-widget
    data-api-url="http://localhost:8080"
    data-username="your-username"
    data-book-type="currently-reading"
    data-max-width="600px"
    data-columns="auto-fill"
    data-min-column-width="100px"
    data-gap="1rem"
    data-show-powered-by="true">
</div>
```

### Available Options

| Option | Default | Description |
|--------|---------|-------------|
| `data-api-url` | Required | Your Hardcover embed server URL |
| `data-username` | Required | Hardcover username to display books for |
| `data-book-type` | `currently-reading` | Book type: 'currently-reading' or 'last-read' |
| `data-max-width` | `800px` | Maximum width of the widget |
| `data-columns` | `auto-fill` | Grid columns (CSS grid value) |
| `data-min-column-width` | `120px` | Minimum width for each book |
| `data-gap` | `1rem` | Space between books |
| `data-show-powered-by` | `true` | Show "Powered by Hardcover" link |

## Examples

### Blog Sidebar

```html
<!-- In your blog sidebar -->
<aside class="sidebar">
    <h3>Currently Reading</h3>
    <div 
        data-hardcover-widget
        data-api-url="http://localhost:8080"
        data-username="your-username"
        data-max-width="300px"
        data-min-column-width="80px">
    </div>
</aside>
<script src="http://localhost:8080/widget.js"></script>
```

### Full Width Footer

```html
<!-- In your footer -->
<footer>
    <div class="container">
        <h2>What I'm Reading</h2>
        <div 
            data-hardcover-widget
            data-api-url="http://localhost:8080"
            data-username="your-username"
            data-columns="6"
            data-gap="2rem">
        </div>
    </div>
</footer>
<script src="http://localhost:8080/widget.js"></script>
```

### Multiple Widgets

You can have multiple widgets on the same page:

```html
<!-- Currently reading books -->
<div 
    data-hardcover-widget
    data-api-url="http://localhost:8080"
    data-username="alice"
    data-book-type="currently-reading"
    data-max-width="400px"
    class="featured-books">
</div>

<!-- Last read books for different user -->
<div 
    data-hardcover-widget
    data-api-url="http://localhost:8080"
    data-username="bob"
    data-book-type="last-read"
    data-max-width="200px"
    data-min-column-width="60px"
    class="sidebar-books">
</div>

<script src="http://localhost:8080/widget.js"></script>
```

### Showing Both Currently Reading and Last Read

Display both types of books for the same user:

```html
<h2>Currently Reading</h2>
<div 
    data-hardcover-widget
    data-api-url="http://localhost:8080"
    data-username="your-username"
    data-book-type="currently-reading">
</div>

<h2>Recently Finished</h2>
<div 
    data-hardcover-widget
    data-api-url="http://localhost:8080"
    data-username="your-username"
    data-book-type="last-read">
</div>

<script src="http://localhost:8080/widget.js"></script>
```

## Styling

The widget comes with default styles, but you can customize it with CSS:

```css
/* Custom widget styles */
.hardcover-widget {
    padding: 2rem;
    background: #f5f5f5;
    border-radius: 12px;
}

/* Customize book covers */
.hw-book-cover {
    border-radius: 8px !important;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1) !important;
}

/* Customize hover effect */
.hw-book-title-overlay {
    background: linear-gradient(to top, #000 0%, transparent 100%) !important;
}
```

## Platform-Specific Instructions

### WordPress

1. Add to a post/page using the HTML block:
```html
<div data-hardcover-widget data-api-url="http://your-server.com" data-username="your-username"></div>
<script src="http://your-server.com/widget.js"></script>
```

2. Or add to your theme's `footer.php`:
```php
<?php if (is_page('about')) : ?>
    <div data-hardcover-widget data-api-url="http://your-server.com" data-username="your-username"></div>
    <script src="http://your-server.com/widget.js"></script>
<?php endif; ?>
```

### React

```jsx
import { useEffect } from 'react';

function CurrentlyReading() {
    useEffect(() => {
        // Load the widget script
        const script = document.createElement('script');
        script.src = 'http://your-server.com/widget.js';
        script.async = true;
        document.body.appendChild(script);

        return () => {
            document.body.removeChild(script);
        };
    }, []);

    return (
        <div 
            data-hardcover-widget
            data-api-url="http://your-server.com"
            data-username="your-username"
            data-max-width="600px"
        />
    );
}
```

### Vue.js

```vue
<template>
    <div 
        data-hardcover-widget
        :data-api-url="apiUrl"
        :data-username="username"
        data-max-width="600px"
    ></div>
</template>

<script>
export default {
    data() {
        return {
            apiUrl: 'http://your-server.com',
            username: 'your-username'
        }
    },
    mounted() {
        const script = document.createElement('script');
        script.src = `${this.apiUrl}/widget.js`;
        script.async = true;
        document.body.appendChild(script);
    }
}
</script>
```

### Static Site Generators

For Hugo, Jekyll, Gatsby, etc., add the HTML to your templates:

```html
<!-- Hugo shortcode example -->
{{< hardcover-widget >}}
    <div data-hardcover-widget data-api-url="{{ .Site.Params.hardcoverApiUrl }}" data-username="{{ .Site.Params.hardcoverUsername }}"></div>
    <script src="{{ .Site.Params.hardcoverApiUrl }}/widget.js"></script>
{{< /hardcover-widget >}}
```

## Troubleshooting

### Books not loading

1. Check that your server is running
2. Verify the `data-api-url` is correct
3. Ensure `data-username` is provided and valid
4. Check browser console for errors
5. Ensure CORS is properly configured on your server

### Styling issues

1. The widget uses CSS custom properties for flexibility
2. Use `!important` if your site's CSS is overriding widget styles
3. Consider using the iframe method for complete style isolation

### Performance

1. The widget lazy-loads images
2. API responses are cached for 30 minutes
3. Consider loading the script with `async` or `defer`

## Security Notes

- The widget makes requests to your Hardcover embed server
- No authentication is required on the client side
- Your API token remains secure on your server
- Enable CORS only for trusted domains in production

## Advanced Usage

### Custom Events

The widget dispatches events you can listen to:

```javascript
document.addEventListener('hardcover:loaded', (e) => {
    console.log('Books loaded:', e.detail.count);
});

document.addEventListener('hardcover:error', (e) => {
    console.error('Loading failed:', e.detail.error);
});
```

### Programmatic Refresh

```javascript
// Get widget instance
const element = document.querySelector('[data-hardcover-widget]');
if (element._hardcoverWidget) {
    element._hardcoverWidget.loadBooks();
}
```

## Production Deployment

When deploying to production:

1. Update all `http://localhost:8080` URLs to your production server URL
2. Configure CORS to only allow your domains
3. Use HTTPS for both the API and widget script
4. Consider using a CDN for the widget.js file
5. Set appropriate cache headers

Example production embed:

```html
<div 
    data-hardcover-widget
    data-api-url="https://books-api.yourdomain.com"
    data-username="your-username">
</div>
<script 
    src="https://books-api.yourdomain.com/widget.js"
    async>
</script>
```
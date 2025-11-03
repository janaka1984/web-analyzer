# Web Analyzer (Gin + GORM + Docker)

An assighment (web analyzer) built in **Go (Golang)** that fetches, analyzes and displays key information from any given website.  
It demonstrates clean architecture concepts, concurrency with goroutines and Go’s native HTTP and HTML parsing features.


## Key Features

- Fetches web pages using Go’s `net/http` client  
- Extracts metadata (doctype, title, headings, links, login forms, etc.)  
- Uses goroutines and channels for concurrent link analysis  
- Displays analysis results via a simple HTML interface (`index.html`, `result.html`)  
- Designed with layered architecture (`adapter`, `service`, `repository`, `transport`)  


## How to Run

### Clone the repo

```bash
git clone https://github.com/janaka1984/web-analyzer.git
cd web-analyzer
```

### Quick start

```bash
docker compose up --build
```

### Open in browser

Visit → http://localhost:8080


## Steps

1 - Enter a website URL on the main page.

2 - Click Analyze.

3 - The app fetches and analyzes the page content, then shows:

- Doctype
- Title
- Headings
- Links (internal/external)
- Presence of login form
- Broken link report

##  Future Improvements

- Web Crawler / Spider Mode:
	Traverse multiple pages automatically and build a site map graph.

- Parallel Link Validation:
	Use worker pools for faster link checking across multiple domains.

- SEO & Accessibility Checks:
	Analyze meta tags, alt text, and Lighthouse-style metrics.

- Dashboard:
	Visualize insights over time.

- API & Frontend Integration:
	Expose REST API endpoints and integrate with React/Vue dashboards.

## Challenges Faced

- Concurrency control: Managing goroutines safely while limiting concurrent requests.

- Network reliability: Handling timeouts and slow responses gracefully.

- HTML parsing consistency: Some sites have malformed HTML or missing metadata.

- SSL/TLS & redirects: Ensuring httpfetch client supports modern HTTPS sites.

- Architecture balance: Keeping the project modular without overcomplicating.


## 🧩 Architectural Diagram

```mermaid
flowchart TD
    subgraph Client
        A["Browser (HTML UI)"]
    end

    subgraph App
        B["Gin Server (Go)"]
        C["Analyzer Service"]
        D["Repository (GORM)"]
    end

    subgraph Infra
        E["PostgreSQL DB"]
        F["Docker Compose Network"]
    end

    A -->|HTTP/HTTPS| B
    B --> C
    C --> D
    D --> E
    B -->|Render HTML| A

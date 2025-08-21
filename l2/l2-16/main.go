package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	baseURL     *url.URL
	maxDepth    int
	visited     map[string]bool
	visitedMux  sync.RWMutex
	client      *http.Client
	userAgent   string
	delay       time.Duration
	downloadDir string
}

type CrawlJob struct {
	URL   string
	Depth int
}

func NewCrawler(baseURL string, maxDepth int, downloadDir string) (*Crawler, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("неверный URL: %v", err)
	}

	// Создаем HTTP клиент с настройками
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Создаем базовую директорию для скачивания, если она не существует
	absDownloadDir, err := filepath.Abs(downloadDir)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения абсолютного пути: %v", err)
	}

	if err := os.MkdirAll(absDownloadDir, 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания базовой директории %s: %v", absDownloadDir, err)
	}

	return &Crawler{
		baseURL:     u,
		maxDepth:    maxDepth,
		visited:     make(map[string]bool),
		client:      client,
		userAgent:   "WebCrawler/1.0 (Go)",
		delay:       time.Second,
		downloadDir: absDownloadDir,
	}, nil
}

func (c *Crawler) isVisited(url string) bool {
	c.visitedMux.RLock()
	defer c.visitedMux.RUnlock()
	return c.visited[url]
}

func (c *Crawler) markVisited(url string) {
	c.visitedMux.Lock()
	defer c.visitedMux.Unlock()
	c.visited[url] = true
}

func (c *Crawler) isSameDomain(targetURL string) bool {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	return u.Host == c.baseURL.Host
}

func (c *Crawler) downloadFile(targetURL string) error {
	// Создаем запрос
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP статус %d для %s", resp.StatusCode, targetURL)
	}

	// Определяем путь для сохранения файла
	filePath, err := c.getFilePath(targetURL)
	if err != nil {
		return fmt.Errorf("ошибка определения пути файла: %v", err)
	}

	// Получаем абсолютный путь для директории
	dir := filepath.Dir(filePath)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("ошибка получения абсолютного пути для %s: %v", dir, err)
	}

	// Создаем директории рекурсивно
	if err := os.MkdirAll(absDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания директории %s: %v", absDir, err)
	}

	// Получаем абсолютный путь для файла
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("ошибка получения абсолютного пути для файла %s: %v", filePath, err)
	}

	// Проверяем, что файл еще не существует
	if _, err := os.Stat(absFilePath); err == nil {
		fmt.Printf("Файл уже существует: %s\n", absFilePath)
		return nil
	}

	// Читаем содержимое
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Если это HTML, обрабатываем ссылки
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "text/html") {
		content := string(body)
		content = c.processHTMLLinks(content, targetURL)
		body = []byte(content)
	}

	// Создаем файл
	file, err := os.Create(absFilePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла %s: %v", absFilePath, err)
	}
	defer file.Close()

	// Записываем содержимое
	_, err = file.Write(body)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	fmt.Printf("Скачан: %s -> %s\n", targetURL, absFilePath)
	return nil
}

func (c *Crawler) getFilePath(targetURL string) (string, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	// Очищаем имя хоста от недопустимых символов для Windows
	host := c.sanitizeFileName(u.Host)

	// Очищаем путь от недопустимых символов и конвертируем / в правильные разделители
	urlPath := strings.Trim(u.Path, "/")
	if urlPath != "" {
		// Заменяем URL-кодированные символы и очищаем каждую часть пути
		pathParts := strings.Split(urlPath, "/")
		for i, part := range pathParts {
			pathParts[i] = c.sanitizeFileName(part)
		}
		urlPath = strings.Join(pathParts, string(filepath.Separator))
	}

	// Создаем базовый путь
	var filePath string
	if urlPath != "" {
		filePath = filepath.Join(c.downloadDir, host, urlPath)
	} else {
		filePath = filepath.Join(c.downloadDir, host)
	}

	// Если путь заканчивается на / или путь пустой, добавляем index.html
	if strings.HasSuffix(u.Path, "/") || u.Path == "" || u.Path == "/" {
		filePath = filepath.Join(filePath, "index.html")
	} else if filepath.Ext(filePath) == "" {
		// Если нет расширения, добавляем .html для HTML страниц
		filePath += ".html"
	}

	return filePath, nil
}

// processHTMLLinks обрабатывает HTML контент и заменяет ссылки на локальные пути
func (c *Crawler) processHTMLLinks(content string, pageURL string) string {
	// Парсим URL текущей страницы
	_, err := url.Parse(pageURL)
	if err != nil {
		return content
	}

	// Регулярные выражения для поиска и замены различных типов ссылок
	patterns := []struct {
		regex   string
		handler func(string, string, string, string) string
	}{
		// href ссылки
		{
			regex:   `(href\s*=\s*["'])([^"']+)(["'])`,
			handler: c.replaceHrefLink,
		},
		// src ссылки для изображений, скриптов и т.д.
		{
			regex:   `(src\s*=\s*["'])([^"']+)(["'])`,
			handler: c.replaceSrcLink,
		},
		// CSS url()
		{
			regex:   `(url\s*\(\s*["']?)([^"')]+)(["']?\s*\))`,
			handler: c.replaceCSSUrl,
		},
	}

	result := content
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern.regex)
		result = re.ReplaceAllStringFunc(result, func(match string) string {
			matches := re.FindStringSubmatch(match)
			if len(matches) >= 4 {
				return pattern.handler(matches[1], matches[2], matches[3], pageURL)
			}
			return match
		})
	}

	return result
}

// replaceHrefLink заменяет href ссылки
func (c *Crawler) replaceHrefLink(prefix, link, suffix, pageURL string) string {
	localPath := c.convertToLocalPath(link, pageURL)
	if localPath != "" {
		return prefix + localPath + suffix
	}
	return prefix + link + suffix
}

// replaceSrcLink заменяет src ссылки
func (c *Crawler) replaceSrcLink(prefix, link, suffix, pageURL string) string {
	localPath := c.convertToLocalPath(link, pageURL)
	if localPath != "" {
		return prefix + localPath + suffix
	}
	return prefix + link + suffix
}

// replaceCSSUrl заменяет CSS url() ссылки
func (c *Crawler) replaceCSSUrl(prefix, link, suffix, pageURL string) string {
	localPath := c.convertToLocalPath(link, pageURL)
	if localPath != "" {
		return prefix + localPath + suffix
	}
	return prefix + link + suffix
}

// convertToLocalPath конвертирует URL в локальный путь
func (c *Crawler) convertToLocalPath(link, pageURL string) string {
	// Пропускаем якоря, mailto, tel и другие специальные ссылки
	if strings.HasPrefix(link, "#") ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "tel:") ||
		strings.HasPrefix(link, "javascript:") ||
		strings.HasPrefix(link, "data:") {
		return ""
	}

	// Разбираем ссылку
	linkURL, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// Разбираем URL текущей страницы
	pageURLParsed, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}

	// Разрешаем относительную ссылку в абсолютную
	absoluteURL := pageURLParsed.ResolveReference(linkURL)

	// Проверяем, что это ссылка на тот же домен
	if absoluteURL.Host != c.baseURL.Host {
		return ""
	}

	// Получаем путь к файлу для этого URL
	targetFilePath, err := c.getFilePath(absoluteURL.String())
	if err != nil {
		return ""
	}

	// Получаем путь к файлу текущей страницы
	currentFilePath, err := c.getFilePath(pageURL)
	if err != nil {
		return ""
	}

	// Вычисляем относительный путь от текущего файла к целевому
	currentDir := filepath.Dir(currentFilePath)
	relativePath, err := filepath.Rel(currentDir, targetFilePath)
	if err != nil {
		return ""
	}

	// Конвертируем Windows пути в веб-формат (/ вместо \)
	relativePath = filepath.ToSlash(relativePath)

	return relativePath
}

func (c *Crawler) extractLinks(content string, baseURL string) []string {
	var links []string

	// Регулярные выражения для поиска ссылок
	patterns := []string{
		`href\s*=\s*["']([^"']+)["']`,
		`src\s*=\s*["']([^"']+)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindAllStringSubmatch(content, -1)

		for _, match := range matches {
			if len(match) > 1 {
				link := match[1]
				absoluteURL := c.resolveURL(link, baseURL)
				if absoluteURL != "" {
					links = append(links, absoluteURL)
				}
			}
		}
	}

	return links
}

func (c *Crawler) resolveURL(link, baseURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// Разрешаем относительные URL
	resolved := base.ResolveReference(u)
	return resolved.String()
}

func (c *Crawler) getPageContent(targetURL string) (string, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	// Читаем только HTML содержимое
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return "", fmt.Errorf("не HTML содержимое: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Crawler) processURL(job CrawlJob) []CrawlJob {
	if job.Depth > c.maxDepth {
		return nil
	}

	if c.isVisited(job.URL) {
		return nil
	}

	if !c.isSameDomain(job.URL) {
		return nil
	}

	c.markVisited(job.URL)

	fmt.Printf("Обрабатываем (глубина %d): %s\n", job.Depth, job.URL)

	// Скачиваем файл
	if err := c.downloadFile(job.URL); err != nil {
		log.Printf("Ошибка скачивания %s: %v", job.URL, err)
		return nil
	}

	// Если достигли максимальной глубины, не извлекаем ссылки
	if job.Depth >= c.maxDepth {
		return nil
	}

	// Получаем содержимое для извлечения ссылок (только для HTML)
	content, err := c.getPageContent(job.URL)
	if err != nil {
		// Не HTML содержимое или ошибка - не извлекаем ссылки
		return nil
	}

	// Извлекаем ссылки
	links := c.extractLinks(content, job.URL)

	var newJobs []CrawlJob
	for _, link := range links {
		if !c.isVisited(link) && c.isSameDomain(link) {
			newJobs = append(newJobs, CrawlJob{
				URL:   link,
				Depth: job.Depth + 1,
			})
		}
	}

	return newJobs
}

func (c *Crawler) Crawl() error {
	jobs := make(chan CrawlJob, 1000) // Увеличил размер буфера
	results := make(chan []CrawlJob, 100)

	// Запускаем воркеры
	const numWorkers = 5
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				newJobs := c.processURL(job)
				if len(newJobs) > 0 {
					select {
					case results <- newJobs:
					case <-time.After(10 * time.Second):
						// Таймаут для предотвращения блокировки
						log.Printf("Воркер %d: таймаут отправки результатов для %s", workerID, job.URL)
					}
				}
				// Задержка между запросами
				time.Sleep(c.delay)
			}
		}(i)
	}

	// Горутина для закрытия канала results после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	// Добавляем начальную задачу
	jobs <- CrawlJob{URL: c.baseURL.String(), Depth: 0}

	// Обрабатываем результаты
	activeJobs := 1
	jobsOpen := true
	processed := 0

	for activeJobs > 0 || jobsOpen {
		select {
		case newJobs, ok := <-results:
			if !ok {
				// Канал results закрыт, все воркеры завершились
				if jobsOpen {
					close(jobs)
					jobsOpen = false
				}
				continue
			}

			for _, job := range newJobs {
				if jobsOpen {
					// Блокирующая отправка - ждем освобождения места в канале
					select {
					case jobs <- job:
						activeJobs++
						processed++
					case <-time.After(30 * time.Second):
						// Только если совсем долго ждем
						log.Printf("Долгое ожидание места в канале для %s, пропускаем", job.URL)
					}
				}
			}
			activeJobs--

			// Логируем прогресс
			if processed%50 == 0 {
				fmt.Printf("Обработано URL: %d, активных задач: %d, размер очереди: %d\n",
					processed, activeJobs, len(jobs))
			}

		case <-time.After(60 * time.Second): // Увеличил таймаут
			// Таймаут для предотвращения бесконечного ожидания
			log.Printf("Таймаут ожидания новых задач (активных: %d, в очереди: %d), завершаем краулинг", activeJobs, len(jobs))
			if jobsOpen {
				close(jobs)
				jobsOpen = false
			}
			// Даем время воркерам завершиться
			time.Sleep(5 * time.Second)
			return nil
		}
	}

	if jobsOpen {
		close(jobs)
	}

	// Ждем завершения всех воркеров
	wg.Wait()

	fmt.Printf("\nВсего обработано URL: %d\n", processed)
	return nil
}

func main() {
	var (
		url         = flag.String("url", "", "URL для скачивания")
		maxDepth    = flag.Int("depth", 2, "Максимальная глубина рекурсии")
		downloadDir = flag.String("dir", "./downloads", "Директория для скачивания")
		help        = flag.Bool("help", false, "Показать помощь")
	)

	flag.Parse()

	if *help || *url == "" {
		fmt.Println("Использование:")
		fmt.Println("  go run main.go -url <URL> [-depth <глубина>] [-dir <директория>]")
		fmt.Println()
		fmt.Println("Параметры:")
		fmt.Println("  -url     URL веб-сайта для скачивания (обязательный)")
		fmt.Println("  -depth   Максимальная глубина рекурсии (по умолчанию: 2)")
		fmt.Println("  -dir     Директория для сохранения файлов (по умолчанию: ./downloads)")
		fmt.Println("  -help    Показать эту справку")
		fmt.Println()
		fmt.Println("Примеры:")
		fmt.Println("  go run main.go -url https://example.com -depth 3")
		fmt.Println("  go run main.go -url https://example.com -depth 1 -dir ./site")
		return
	}

	// Создаем краулер
	crawler, err := NewCrawler(*url, *maxDepth, *downloadDir)
	if err != nil {
		log.Fatalf("Ошибка создания краулера: %v", err)
	}

	fmt.Printf("Начинаем скачивание сайта: %s\n", *url)
	fmt.Printf("Максимальная глубина: %d\n", *maxDepth)
	fmt.Printf("Директория сохранения: %s\n", *downloadDir)
	fmt.Println("Нажмите Ctrl+C для остановки")
	fmt.Println()

	// Запускаем краулинг
	if err := crawler.Crawl(); err != nil {
		log.Fatalf("Ошибка краулинга: %v", err)
	}

	fmt.Println("\nСкачивание завершено!")
}

// sanitizeFileName очищает строку от недопустимых символов для файловой системы
func (c *Crawler) sanitizeFileName(name string) string {
	// Разрешаем только буквы, цифры, точку, дефис и подчёркивание
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	clean := re.ReplaceAllString(name, "_")

	// Ограничим длину имени, чтобы не было слишком длинных путей
	if len(clean) > 200 {
		clean = clean[:200]
	}

	// На всякий случай удалим пробелы по краям
	return strings.TrimSpace(clean)
}

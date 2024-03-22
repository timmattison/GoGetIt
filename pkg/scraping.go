package lib

import (
	"errors"
	"github.com/playwright-community/playwright-go"
	"os"
	"strings"
	"time"
)

var (
	ErrFailedToClickElement     = errors.New("failed to click element")
	ErrCouldNotGetContent       = errors.New("could not get content")
	ErrCouldNotGoBack           = errors.New("could not go back")
	ErrContentDidNotLoad        = errors.New("content did not load")
	ErrCouldNotGetAttribute     = errors.New("could not get attribute")
	ErrCouldNotGetInnerText     = errors.New("could not get inner text")
	ErrCouldNotCountElements    = errors.New("could not count elements")
	ErrMoreThanOneElement       = errors.New("expected one element, got more than one")
	ErrCouldNotCreateScreenshot = errors.New("could not create screenshot")
)

func ClickAndWait(page playwright.Page, locator playwright.Locator, loadingTimeout time.Duration) error {
	var err error

	if err = locator.Click(); err != nil {
		return errors.Join(ErrFailedToClickElement, err)
	}

	return WaitForPageToLoad(page, loadingTimeout)
}

func GetContent(page playwright.Page) (string, error) {
	var content string
	var err error

	if content, err = page.Content(); err != nil {
		return "", errors.Join(ErrCouldNotGetContent, err)
	}

	return content, nil
}

func GoBackAndWait(page playwright.Page, loadingTimeout time.Duration) error {
	var err error

	if _, err = page.GoBack(); err != nil {
		return errors.Join(ErrCouldNotGoBack, err)
	}

	return WaitForPageToLoad(page, loadingTimeout)
}

func WaitForPageToLoad(page playwright.Page, loadingTimeout time.Duration) error {
	var err error

	if err = waitForLoadState(page, loadingTimeout); err != nil {
		return errors.Join(ErrContentDidNotLoad, err)
	}

	return nil
}

func GetAttribute(locator playwright.Locator, attributeName string) (string, error) {
	var value string
	var err error

	if value, err = locator.GetAttribute(attributeName); err != nil {
		return "", errors.Join(ErrCouldNotGetAttribute, err)
	}

	return value, nil
}

func GetInnerText(locator playwright.Locator) (string, error) {
	var value string
	var err error

	if value, err = locator.InnerText(); err != nil {
		return "", errors.Join(ErrCouldNotGetInnerText, err)
	}

	return strings.TrimSpace(value), nil
}

func FilterBySelector(page playwright.Page, selector string, filterFunc func(playwright.Locator) bool) ([]playwright.Locator, error) {
	elements := page.Locator(selector)

	// Get the count of elements
	var elementCount int
	var err error

	if elementCount, err = elements.Count(); err != nil {
		return nil, errors.Join(ErrCouldNotCountElements, err)
	}

	var result []playwright.Locator

	for i := range elementCount {
		element := elements.Nth(i)

		if filterFunc != nil && !filterFunc(element) {
			continue
		}

		result = append(result, element)
	}

	return result, nil
}

func FilterBySelectorExactlyOne(page playwright.Page, selector string, filterFunc func(playwright.Locator) bool) (playwright.Locator, error) {
	var result []playwright.Locator
	var err error

	if result, err = FilterBySelector(page, selector, filterFunc); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, ErrMoreThanOneElement
	}

	return result[0], nil
}

func Screenshot(page playwright.Page, filename string) error {
	var err error

	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(filename),
	}); err != nil {
		return ErrCouldNotCreateScreenshot
	}

	return nil
}

func waitForLoadState(page playwright.Page, timeout time.Duration) error {
	return page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
}

func SaveContent(page playwright.Page, filename string) error {
	var content string
	var err error

	if content, err = GetContent(page); err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

const lunr = require('lunr');

const searchInput = document.getElementById('search-input');
const searchResultsMessage = document.getElementById('search-results-message');
const searchResultsCount = document.getElementById('search-results-count');
const searchResultsList = document.getElementById('search-results');
const searchResults = Array.from(searchResultsList.children)

let searchIndex;
let results;

function showMessage(message) {
    if (!message) {
        searchResultsMessage.classList.add('hidden');
        return;
    }

    searchResultsMessage.textContent = message;
    searchResultsMessage.classList.remove('hidden');

    searchResultsCount.classList.add('hidden');
}

function getInitialSearchTerm() {
    const params = new URLSearchParams(window.location.search);
    return params.get('q') || '';
}

function onSearch(term) {
    console.log('Searching for', term);
    searchResults.forEach(elm => elm.classList.add('hidden'));
    if (term.length < 3) {
        showMessage('Search term must be at least 3 characters');
        return;
    }
    showMessage('Searching...');

    results = searchIndex.search(term);

    if (!results.length) {
        showMessage('No results found');
        return;
    }

    console.log(results);

    searchResultsMessage.classList.add('hidden');
    searchResultsCount.textContent = `${results.length} result${results.length === 1 ? '' : 's'} found`;
    searchResultsCount.classList.remove('hidden');

    searchResults.forEach(function (elm) {
        if (results.filter(res => res.ref === elm.dataset.id).length > 0) {
            elm.classList.remove('hidden');
            elm.classList.add('block');
        } else {
            elm.classList.add('hidden');
            elm.classList.remove('block');
        }
    });

    let sortOrder = Array.from(document.getElementsByName('sort-by')).filter(elm => elm.checked)[0].id.split('-')[1];
    searchResults.sort(function (b, a) {
        switch(sortOrder) {
            case "cost":
                return a.dataset.cost - b.dataset.cost;
            case "date":
            default:
                return new Date(a.dataset.date) - new Date(b.dataset.date);
        }
    });

    searchResults.forEach(function (result) {
        searchResultsList.appendChild(result);
    });

    window.history.replaceState({}, '', `?q=${term}`);
}

searchResults.forEach(function (elm) {
    elm.classList.add('hidden');
    elm.classList.remove('block');
});

window.addEventListener('load', function () {
    document.getElementById('search-form').addEventListener('submit', function (event) {
        event.preventDefault();
    });

    document.getElementsByName('sort-by').forEach(function (elm) {
        elm.addEventListener('change', function () {
            onSearch(searchInput.value);
        });
    });

    searchInput.disabled = true;
    const initialSearchterm = getInitialSearchTerm();
    searchInput.value = initialSearchterm;

    fetch('/search/index.json')
        .then(response => response.json())
        .then(data => {
            searchIndex = lunr(function () {
                this.ref('trip');
                this.field('justification');
                this.field('passengers');
                this.field('route');

                data.forEach(res => {
                    this.add(res);
                }, this);
            });

            searchInput.addEventListener('input', function () {
                onSearch(this.value);
            });

            searchInput.disabled = false;
            searchInput.focus();

            onSearch(searchInput.value);
        });
});

const form = document.getElementById("searchForm");
const results = document.getElementById("results");
const genreInput = document.getElementById("genres");
const addGenreButton = document.getElementById("addGenre");
const selectedGenresContainer = document.getElementById("selectedGenres");

let selectedGenres = [];

genreInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        event.preventDefault();
        addGenre();
    }
});

addGenreButton.addEventListener("click", addGenre);

function addGenre() {
    const genre = genreInput.value.trim();

    if (!genre) return;

    // Prevent duplicates
    if (selectedGenres.some(g => g.toLowerCase() === genre.toLowerCase())) {
        return;
    }

    selectedGenres.push(genre);

    genreInput.value = "";
    renderGenres();
}

// Display selected genres
function renderGenres() {
    selectedGenresContainer.replaceChildren();

    selectedGenres.forEach(genre => {
        const chip = document.createElement("div");
        chip.className = "genre-chip";

        const label = document.createElement("span");
        label.textContent = genre;

        const removeButton = document.createElement("button");
        removeButton.type = "button";
        removeButton.textContent = "×";

        removeButton.addEventListener("click", () => {
            selectedGenres = selectedGenres.filter(g => g !== genre);
            renderGenres();
        });

        chip.append(label, removeButton);
        selectedGenresContainer.appendChild(chip);
    });
}

// Search for movies
form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const origin = document.getElementById("origin").value.trim();
    const director = document.getElementById("director").value.trim();

    var directorId = 0;
    if (director) {
        try {
        const response = await fetch(`/director?with_crew=${director}`);
        if (!response.ok) {
            throw new Error(await response.text());
        }
        const data = await response.json();
        directorId = data.results[0].id
        } catch (error) {
        results.textContent = "Error: " + error.message;
        }
    }
    
    const params = new URLSearchParams();

    if (selectedGenres.length > 0) {
        params.set("with_genres", selectedGenres.join(","));
    }

    if (origin) {
        params.set("with_origin_country", origin);
    }

    if (directorId) {
        params.set("with_crew", directorId)
    }

    results.textContent = "Loading...";

    try {
        const response = await fetch(`/recs?${params.toString()}`);

        if (!response.ok) {
            throw new Error(await response.text());
        }

        const data = await response.json();
        renderMovies(data.results);

    } catch (error) {
        results.textContent = "Error: " + error.message;
    }
});

function renderMovies(movies) {
    results.replaceChildren();
    if (movies.length === 0) {
        results.textContent = "No movies found.";
        return;
    }

    movies.forEach(movie => {
        // Create the movie container
        const movieElement = document.createElement("div");
        movieElement.className = "movie";

        // Create the poster
        const poster = document.createElement("img");

        if (movie.poster_path) {
            poster.src = `https://image.tmdb.org/t/p/w185${movie.poster_path}`;
            poster.alt = `Poster for ${movie.title}`;
        } else {
            poster.alt = "No poster available";
        }

        // Create the movie information
        const info = document.createElement("div");

        const title = document.createElement("h3");
        title.textContent = movie.title;

        const releaseDate = document.createElement("p");
        releaseDate.textContent = `Release date: ${movie.release_date || "Unknown"}`;

        const rating = document.createElement("p");
        rating.textContent = `Rating: ${movie.vote_average}/10`;

        const synopsis = document.createElement("p");
        synopsis.textContent = `${movie.overview}`

        // Put everything together
        info.append(title, releaseDate, rating, synopsis);
        movieElement.append(poster, info);

        // Add the movie to the results container
        results.appendChild(movieElement);
    });
}

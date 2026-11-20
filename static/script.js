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
        removeButton.setAttribute("aria-label", `Remove ${genre}`);

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
    const releaseYear = document.getElementById("year").value;
    const actor = document.getElementById("actor").value.trim();
    const maxRuntime = document.getElementById("maxRuntime").value;
    const minRuntime = document.getElementById("minRuntime").value;

    const params = new URLSearchParams();

    if (director) {
        var directorID = await fetchPersonID(director);
        params.set("with_crew", directorID)
    }

    if (actor) {
        var actorID = await fetchPersonID(actor);
        params.set("with_cast", actorID)
    }

    if (selectedGenres.length > 0) {
        params.set("with_genres", selectedGenres.join(","));
    }

    if (origin) {
        params.set("with_origin_country", origin);
    }

    if (releaseYear) {
        params.set("primary_release_year", releaseYear)
    }

    if (maxRuntime) {
        params.set("maxRuntime", maxRuntime)
    }

    if (minRuntime) {
        params.set("minRuntime", minRuntime)
    }

    console.log("Parameters:", params.toString());
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

// Builds <p><strong>Label</strong> value</p>
function makeLine(label, value) {
    const p = document.createElement("p");
    const strong = document.createElement("strong");
    strong.textContent = label + " ";
    p.append(strong, value);
    return p;
}

function renderMovies(movies) {
    results.replaceChildren();
    if (!movies || movies.length === 0) {
        results.textContent = "No movies found.";
        return;
    }

    movies.forEach(movie => {
        // Create the movie card
        const movieElement = document.createElement("div");
        movieElement.className = "movie";

        // Poster
        const poster = document.createElement("img");

        if (movie.poster_path) {
            poster.src = `https://image.tmdb.org/t/p/w185${movie.poster_path}`;
            poster.alt = `Poster for ${movie.title}`;
        } else {
            poster.alt = "No poster available";
        }

        // Details column
        const info = document.createElement("div");
        info.className = "movie-info";

        const title = document.createElement("h3");
        title.textContent = movie.title;

        const releaseDate = makeLine("Released:", movie.release_date || "Unknown");

        const rating = makeLine("Rating:", `${movie.vote_average}/10`);

        const directors = (movie.crew ?? []).map(person => person.name);
        const crew = makeLine("Directed by:", directors.join(", ") || "Unknown");

        const actors = (movie.cast ?? []).map(person => person.name);
        const cast = makeLine("Starring:", actors.join(", ") || "Unknown");

        info.append(title, releaseDate, rating, crew, cast);

        // Runtime is 0 when TMDB has no value, so only show real ones
        if (movie.runtime > 0) {
            info.appendChild(makeLine("Runtime:", `${movie.runtime} minutes`));
        }

        const genres = document.createElement("p");
        genres.className = "genres";
        genres.textContent = (movie.genres ?? []).join(", ");
        info.appendChild(genres);

        // Synopsis is its own column, so it goes next to info, not inside it
        const synopsis = document.createElement("p");
        synopsis.className = "movie-synopsis";
        synopsis.textContent = movie.overview || "No synopsis available.";

        movieElement.append(poster, info, synopsis);

        // Add the movie to the results container
        results.appendChild(movieElement);
    });
}

async function fetchPersonID(name) {
    var personID = 0;
    try {
        const response = await fetch(`/person?person=${encodeURIComponent(name)}`);
        if (!response.ok) {
            throw new Error(await response.text());
        }
        const data = await response.json();
        personID = data.results[0].id
    } catch (error) {
        results.textContent = "Error: " + error.message;
    }
    return personID;
}

class BookService {
  urlApi;

  constructor() {
    this.urlApi = process.env["BACKEND_URL"] || "backend"; //"http://127.0.0.1:8080"
  }

  searchBook(word:string) {
    return fetch(`${this.urlApi}/Search/${word}`, {
      mode: "cors",
    })
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }

  searchBookRegex(regex:string) {
    return fetch(`${this.urlApi}/RegexSearch/${regex}`, {
      mode: "cors",
    })
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }
}

export default BookService;

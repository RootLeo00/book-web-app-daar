class BookService {
  urlApi;

  constructor() {
    this.urlApi = process.env["BACKEND_URL"] || "backend"; //"http://127.0.0.1:8000"
  }

  searchBook(word:string) {
    return fetch(`${this.urlApi}/Search/${word}/`)
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }

  searchBookRegex(regex:string) {
    return fetch(`${this.urlApi}/RegexSearch/${regex}/`)
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }
}

export default BookService;
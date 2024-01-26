import React, { useEffect, useState } from 'react';
import { IBook } from '../models/book';
import BookService from '../book.service';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import '../styles/styles.css';
import Slider from '@mui/material/Slider';
import FormControlLabel from '@mui/material/FormControlLabel';
import Stack from '@mui/material/Stack';
import BookCard from '../components/BookCard';
import { FormControl, Radio, RadioGroup } from '@mui/material';

const HomePage = () => {
    const [data, setData] = useState<IBook[]>([]);
    const [inputNameBook, setInputNameBook] = useState<string>("");
    const [inputRegex, setInputRegex] = useState<string>("");
    const [suggestions, setSuggestions] = useState<IBook[]>([]);
    const [numberOfBooks, setNumberOfBooks] = useState<number>(5);
    const [sort, setSort] = useState<string>("none");
    const [maxNumberOfBooks, setMaxNumberOfBooks] = useState<number>(20);

    const bookService = new BookService();

    const getbooks = async (word: string, regex: boolean) => {
        try {
            if (regex) {
                console.log("searching for regex: ", word)
                bookService.searchBookRegex(word)
                    .then(async (res) => {
                        const obj = await JSON.parse(JSON.stringify(res))
                        setData(obj.books);
                        setSuggestions(obj.neighbors);
                        setMaxNumberOfBooks(obj.books.length);
                    })
            }
            else {
                console.log("searching for name: ", word)
                bookService.searchBook(word)
                    .then(async (res) => {
                        console.log(res)
                        
                        const obj = await JSON.parse(JSON.stringify(res))

                        setData(obj.books);
                        setSuggestions(obj.neighbors);
                        setMaxNumberOfBooks(obj.books.length);
                    })
            }

        } catch (err) {
            alert('failed loading json data');
            console.log(err);
        }
    };

    useEffect(() => {
        // Sorting the data based on occurrence or pertinence when the corresponding switch is true
        const updatedSortedData = [...data].sort((a, b) => {
            if (sort === "pertinence") {
                return b.crank - a.crank;
            } else if (sort === "occurrence") {
                return b.occurrence - a.occurrence;
            }
            return 0; // No sorting 
        });

        setData(updatedSortedData);

        const updatedSortedSuggestions = (suggestions !== undefined) ? ([...suggestions].sort((a, b) => {
            if (sort === "pertinence") {
                return b.crank - a.crank;
            } else if (sort === "occurrence") {
                return b.occurrence - a.occurrence;
            }
            return 0; // No sorting 
        })) : ([]);

        setSuggestions(updatedSortedSuggestions);

    }, [suggestions, data, sort]);


    return (
        <div className="search-section">
            <div style={{ display: 'flex', alignItems: 'center' }}>
                <TextField
                    type="text"
                    id="namebook_res"
                    label="Name of the Book to Search"
                    color='warning'
                    value={inputNameBook}
                    onChange={(e) => { setInputNameBook(e.target.value) }}
                />
                <Button onClick={() => getbooks(inputNameBook, false)}>Search Book Name</Button>
            </div>

            <div style={{ display: 'flex', alignItems: 'center' }}>
                <TextField
                    type="text"
                    id="regex_res"
                    label="String to Search in all of the Books"
                    color='warning'
                    value={inputRegex}
                    onChange={(e) => { setInputRegex(e.target.value) }}
                />
                <Button onClick={() => getbooks(inputRegex, true)}>Search String</Button>
            </div>
            <br />
            <br />

            <div style={{ display: 'flex', alignItems: 'center' }}>
                <FormControlLabel
                    control={
                        <Slider
                            value={numberOfBooks}
                            onChange={(e, value) => { setNumberOfBooks(value as number) }} // TODO: make it more elegant
                            min={0}
                            max={maxNumberOfBooks}
                            valueLabelDisplay="auto"
                            aria-labelledby="numberOfBooks-slider"
                            color="warning"
                        />
                    }
                    label="Number of books to display"
                />
                <FormControl>
                    <RadioGroup
                        aria-label="sortingOption"
                        name="sortingOption"
                        value={sort}
                        onChange={(event,value) => { setSort(value) }}
                        row
                    >
                        <FormControlLabel
                            value="occurrence"
                            control={<Radio color="warning" />}
                            label="Sort by Occurrence"
                        />
                        <FormControlLabel
                            value="pertinence"
                            control={<Radio color="warning" />}
                            label="Sort by Pertinence"
                        />
                    </RadioGroup>
                </FormControl>
            </div>
            <br />
            <br />
            <div>
                <h2>Books</h2>
                <Stack spacing={{ xs: 1, sm: 2 }} direction="row" useFlexGap flexWrap="wrap">
                    {data.slice(0, numberOfBooks).map((book, index) => (
                        <BookCard key={index} cardData={book} />
                    ))}
                </Stack>
            </div>
            <div>
                <h2>Suggested Books</h2>
                <Stack spacing={{ xs: 1, sm: 2 }} direction="row" useFlexGap flexWrap="wrap">
                    {suggestions?.slice(0, numberOfBooks).map((book, index) => (
                        <BookCard key={index} cardData={book} />
                    ))}
                </Stack>
            </div>
        </div >
    );
};

export default HomePage;

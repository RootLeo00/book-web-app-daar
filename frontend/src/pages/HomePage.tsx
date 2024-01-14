import React, { useState, useEffect } from 'react';
import { IBook } from '../models/book';
import BookService from '../book.service';
import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import TextField from '@mui/material/TextField';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import '../styles/styles.css';

const HomePage = () => {
    const [data, setData] = useState<IBook[]>([]);
    const [dataS, setDataS] = useState<IBook[]>([]);
    const [suggestions, setSuggestions] = useState<IBook[]>([]);
    const [visi, setVisi] = useState(-1);

    const bookService = new BookService();

    const getbooks = async (word: string) => {
        try {
            setData([]);
            setSuggestions([]);

            const res = await bookService.searchBook(word);
            const obj = JSON.parse(JSON.stringify(res));
            setData(obj.books);
            setDataS(obj.books);
            setSuggestions(obj.neightboors);
            console.log(suggestions);
        } catch (err) {
            alert('failed loading json data');
            console.log(err);
        }
    };

    const handleClick = () => {
        const isChecked = (document.getElementById("toggle_checkbox") as HTMLInputElement).checked;
        console.log(isChecked);

        setData([]);
        setSuggestions([]);

        if (isChecked) {
            setVisi(1);
        } else {
            setVisi(-1);
        }
    };

    const getbooksR = async (regex: string) => {
        try {
            const res = await bookService.searchBookRegex(regex);
            const obj = JSON.parse(JSON.stringify(res));
            setData(obj.books);
            setDataS(obj.books);
            setSuggestions(obj.neightboors);
            console.log(suggestions);
        } catch (err) {
            alert('failed loading json data');
            console.log(err);
        }
    };

    const numberOfBooksRes = () => {
        const e = (document.getElementById("numberbook_res") as HTMLInputElement).value;

        if (e !== "max") {
            setData(dataS.slice(0, parseInt(e, 10)));
        } else {
            setData(dataS);
        }

        (document.getElementById("numberbook_res") as HTMLInputElement).value = "def";
    };

    const sortbyOccurrenceRes = () => {
        const e = (document.getElementById("occ_res") as HTMLInputElement).value;

        if (e === "ON") {
            setData([...data].sort((book1, book2) => book2.occurrence - book1.occurrence));
        } else {
            setData(dataS);
        }

        (document.getElementById("occ_res") as HTMLInputElement).value = "def";
    };

    const sortbyPertinenceRes = () => {
        const e = (document.getElementById("pert_res") as HTMLInputElement).value;

        if (e === "ON") {
            setData([...data].sort((book1, book2) => book2.crank - book1.crank));
        } else {
            setData(dataS);
        }

        (document.getElementById("pert_res") as HTMLInputElement).value = "def";
    };

    useEffect(() => {
        // componentDidMount logic here
    }, []);

    return (
        <div className="search-section">
            <Button onClick={() => getbooks("yourSearchWord")}>Search Books</Button>

            <FormControl>
                <InputLabel htmlFor="toggle_checkbox">Suggestions</InputLabel>
                <Checkbox
                    id="toggle_checkbox"
                    onChange={handleClick}
                />
            </FormControl>

            <Button onClick={() => getbooksR("yourRegex")}>Search Books with Regex</Button>

            <TextField
                type="text"
                id="numberbook_res"
                label="Number of Books"
            />
            <Button onClick={numberOfBooksRes}>Apply Number of Books</Button>

            <FormControl>
                <InputLabel htmlFor="occ_res">Sort by Occurrence</InputLabel>
                <Select id="occ_res" defaultValue="OFF">
                    <MenuItem value="OFF">Off</MenuItem>
                    <MenuItem value="ON">On</MenuItem>
                </Select>
            </FormControl>
            <Button onClick={sortbyOccurrenceRes}>Apply Sort by Occurrence</Button>

            <FormControl>
                <InputLabel htmlFor="pert_res">Sort by Pertinence</InputLabel>
                <Select id="pert_res" defaultValue="OFF">
                    <MenuItem value="OFF">Off</MenuItem>
                    <MenuItem value="ON">On</MenuItem>
                </Select>
            </FormControl>
            <Button onClick={sortbyPertinenceRes}>Apply Sort by Pertinence</Button>

            {/* Your JSX goes here */}
        </div>
    );
};

export default HomePage;

import { TextField } from "@mui/material";

function Search({ id = "Search", label = "Search", onChange }) {
    return (
        <TextField
            id={id}
            sx={{
                display: "flex",
                width: "calc(100% - 1rem)",
                margin: ".5rem",
                alignSelf: "center",
            }}
            label={label}
            type="search"
            size="small"
            //value={value}
            onChange={(e) => onChange(e.target.value)}
        />
    );
}

export default Search;

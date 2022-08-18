import "./SearchBar.scss";

import { BsSearch } from "react-icons/bs";

const SearchBar = ({ filter }: { filter: string }) => {
  return (
    <div className="searchBar">
      <BsSearch className="icon" />
      <input placeholder={`Filter ${filter}`} />
    </div>
  );
};

export default SearchBar;

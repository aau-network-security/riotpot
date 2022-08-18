import SearchBar from "../../components/searchbar/SearchBar";
import ServicesHeader from "./ServicesHeader";
import { ServicesTable } from "./ServicesTable";
import { ServicesUtils } from "./ServicesUtils";

import "./Services.scss";

const Services = () => {
  const view: string = "Services";

  return (
    <main>
      <ServicesHeader view={view} />
      <ServicesUtils view={view} />
      <SearchBar filter="services and proxies" />
      <ServicesTable />
    </main>
  );
};

export default Services;

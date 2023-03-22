import SearchBar from "../../components/searchbar/SearchBar";
import ServicesHeader from "./ServicesHeader";
import { ServicesTable } from "./ServicesTable";
import { ServicesUtils } from "./ServicesUtils";

import "./Services.scss";
import { useRecoilValue } from "recoil";
import { services } from "../../recoil/atoms/services";

const Services = () => {
  const servicesList = useRecoilValue(services);
  const view: string = "Services";

  return (
    <main>
      <ServicesHeader view={view} />
      <ServicesUtils />
      <SearchBar filter="services and proxies" />
      <ServicesTable servicesList={servicesList} />
    </main>
  );
};

export default Services;

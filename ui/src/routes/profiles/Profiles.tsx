import SearchBar from "../../components/searchbar/SearchBar";
import Title from "../../components/title/Title";
import { ProfilesTable } from "./ProfilesTable";
import { ProfilesUtils } from "./ProfilesUtils";

const Profiles = () => {
  // Title and subtitle
  const title: string = "Profiles";
  const subTitle: string = "List of profiles available";

  return (
    <main>
      <Title title={title} subTitle={subTitle} />
      <ProfilesUtils />
      <SearchBar filter="instances" />
      <ProfilesTable />
    </main>
  );
};

export default Profiles;

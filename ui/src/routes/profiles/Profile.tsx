import { useParams } from "react-router-dom";
import { useRecoilValue } from "recoil";
import Title from "../../components/title/Title";
import { profilesFilter } from "../../recoil/atoms/profiles";
import { ServicesTable } from "../services/ServicesTable";
import { ProfileUtils } from "./ProfileUtils";

const Profile = () => {
  let { id } = useParams();

  const profile = useRecoilValue(profilesFilter(id));

  return (
    <main>
      <Title title={profile.name} subTitle={profile.description} />
      <ProfileUtils id={profile.id} />
      <ServicesTable servicesList={profile.services} />
    </main>
  );
};

export default Profile;

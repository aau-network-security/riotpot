import { useParams } from "react-router-dom";
import { useRecoilValue } from "recoil";
import Title from "../../components/title/Title";
import { instances } from "../../recoil/atoms/instances";

const Instance = () => {
  let { id } = useParams();

  const idN = id ? parseInt(id) : -1;
  const instance = useRecoilValue(instances(idN));

  return (
    <main>
      <Title title={instance.name} subTitle={instance.description} />
    </main>
  );
};

export default Instance;

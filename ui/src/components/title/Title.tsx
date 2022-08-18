import "./Title.scss";

type TitleProps = {
  title: string;
  subTitle?: string;
};

const Title = ({ title, subTitle }: TitleProps) => {
  return (
    <div className="pageTitle">
      <h3 className="title">{title}</h3>
      <p className="subTitle">{subTitle}</p>
    </div>
  );
};

export default Title;

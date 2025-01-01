import { Spin } from 'antd';

export const Loading = () => {
  return (
    <div className="flex justify-center items-center h-full">
      <Spin size="large" />
    </div>
  );
};
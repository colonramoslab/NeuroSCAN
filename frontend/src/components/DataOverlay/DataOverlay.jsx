import React from 'react';
import { useSelector } from 'react-redux';

const DataOverlay = () => {
  const data = useSelector((state) => state.dataOverlay);
  console.log({ data });

  const title = 'honk';
  return (
    <h1>{title}</h1>
  );
};

export default DataOverlay;

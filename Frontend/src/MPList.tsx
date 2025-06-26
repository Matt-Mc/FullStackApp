import React from 'react';
import type { MP } from './types';

interface MPListProps {
  mps: MP[];
}

const MPList: React.FC<MPListProps> = ({ mps }) => {
  return (
    <div className="mp-list">
      {mps.map((mp) => (
        <div key={mp.name} className="mp-item">
          <img src={mp.image} alt={mp.name} />
          <h2>{mp.name}</h2>
          <p>{mp.current_party.short_name.en} - {mp.current_party.short_name.en}, {mp.current_riding.province}</p>
          <p><a href={mp.url} target="_blank" rel="noopener noreferrer">View Profile</a></p>
        </div>
      ))}
    </div>
  );
};

export default MPList;

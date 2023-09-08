'use client';
import SubPageHeader from '@/components/headers/SubPageHeader';
import { PageProps } from '@/components/types';
import { ReactElement } from 'react';
import DestinationConnectionCard from './components/DestinationConnectionCard';

export default function Page({ params }: PageProps): ReactElement {
  const id = params?.id ?? '';
  return (
    <div className="job-details-container">
      <SubPageHeader
        header="Destination Connections"
        description={`Manage job's destination connections`}
      />

      <div className="space-y-10">
        <DestinationConnectionCard jobId={id} />
      </div>
    </div>
  );
}

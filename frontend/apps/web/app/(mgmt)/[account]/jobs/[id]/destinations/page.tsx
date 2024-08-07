'use client';
import ButtonText from '@/components/ButtonText';
import SubPageHeader from '@/components/headers/SubPageHeader';
import { useAccount } from '@/components/providers/account-provider';
import { PageProps } from '@/components/types';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useQuery } from '@connectrpc/connect-query';
import { getConnections, getJob } from '@neosync/sdk/connectquery';
import { PlusIcon } from '@radix-ui/react-icons';
import NextLink from 'next/link';
import { ReactElement } from 'react';
import {
  getConnectionIdFromSource,
  getFkIdFromGenerateSource,
} from '../source/components/util';
import { isAiDataGenJob, isDataGenJob } from '../util';
import DestinationConnectionCard from './components/DestinationConnectionCard';

export default function Page({ params }: PageProps): ReactElement {
  const id = params?.id ?? '';
  const { account } = useAccount();
  const {
    data,
    isLoading,
    refetch: mutate,
  } = useQuery(getJob, { id }, { enabled: !!id });
  const { data: connectionsData, isLoading: isConnectionsLoading } = useQuery(
    getConnections,
    { accountId: account?.id },
    { enabled: !!account?.id }
  );

  const connections = connectionsData?.connections ?? [];
  const destinationIds = new Set(
    data?.job?.destinations.map((d) => d.connectionId)
  );
  const sourceConnectionId = getConnectionIdFromSource(data?.job?.source);
  const fkConnectionId = getFkIdFromGenerateSource(data?.job?.source);
  const fkConnection = connections.find((c) => c.id === fkConnectionId);
  return (
    <div className="job-details-container">
      <SubPageHeader
        header="Destination Connections"
        description={`Manage a job's destination connections`}
        extraHeading={
          isDataGenJob(data?.job) || isAiDataGenJob(data?.job) ? null : (
            <NewDestinationButton jobId={id} />
          )
        }
      />

      {isLoading || isConnectionsLoading ? (
        <Skeleton className="w-full h-96 rounded-lg" />
      ) : (
        <div className="space-y-10">
          {data?.job?.destinations.map((destination) => {
            return (
              <DestinationConnectionCard
                key={destination.id}
                jobSourceId={
                  fkConnectionId
                    ? fkConnectionId
                    : sourceConnectionId
                      ? sourceConnectionId
                      : ''
                }
                jobId={id}
                destination={destination}
                mutate={mutate}
                connections={connections}
                availableConnections={connections.filter((c) => {
                  if (isDataGenJob(data?.job) || isAiDataGenJob(data?.job)) {
                    return (
                      c.connectionConfig?.config.case ===
                      fkConnection?.connectionConfig?.config.case
                    );
                  }
                  return (
                    c.id === destination.connectionId ||
                    (c.id != sourceConnectionId && !destinationIds?.has(c.id))
                  );
                })}
                isDeleteDisabled={data?.job?.destinations.length === 1}
              />
            );
          })}
        </div>
      )}
    </div>
  );
}

interface NewDestinationButtonProps {
  jobId: string;
}

function NewDestinationButton(props: NewDestinationButtonProps): ReactElement {
  const { jobId } = props;
  const { account } = useAccount();
  return (
    <NextLink href={`/${account?.name}/new/job/${jobId}/destination`}>
      <Button>
        <ButtonText leftIcon={<PlusIcon />} text="New Destination" />
      </Button>
    </NextLink>
  );
}

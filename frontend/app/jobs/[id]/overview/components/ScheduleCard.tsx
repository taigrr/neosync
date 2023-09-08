'use client';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { useToast } from '@/components/ui/use-toast';
import {
  UpdateJobScheduleRequest,
  UpdateJobScheduleResponse,
} from '@/neosync-api-client/mgmt/v1alpha1/job_pb';
import { getErrorMessage } from '@/util/util';
import { yupResolver } from '@hookform/resolvers/yup';
import cron from 'cron-validate';
import { ReactElement } from 'react';
import { useForm } from 'react-hook-form';
import * as Yup from 'yup';
import { getJob } from '../util';

const SCHEDULE_FORM_SCHEMA = Yup.object({
  cronSchedule: Yup.string()
    .optional()
    .test('isValidCron', 'Not a valid cron schedule', (value) => {
      if (!value) {
        return true;
      }
      return !!value && cron(value).isValid();
    }),
});

export type ScheduleFormValues = Yup.InferType<typeof SCHEDULE_FORM_SCHEMA>;

interface Props {
  jobId: string;
}

export default function JobScheduleCard({ jobId }: Props): ReactElement {
  const { toast } = useToast();
  const form = useForm({
    resolver: yupResolver<ScheduleFormValues>(SCHEDULE_FORM_SCHEMA),
    defaultValues: async () => {
      const res = await getJob(jobId);
      if (!res) {
        return { cronSchedule: '' };
      }
      return { cronSchedule: res.job?.cronSchedule || '' };
    },
  });

  async function onSubmit(values: ScheduleFormValues) {
    try {
      await updateJobSchedule(jobId, values.cronSchedule);
      toast({
        title: 'Successfully updated job schedule!',
        variant: 'default',
      });
    } catch (err) {
      console.error(err);
      toast({
        title: 'Unable to update job schedule',
        description: getErrorMessage(err),
        variant: 'destructive',
      });
    }
  }

  const msg =
    !form.getValues().cronSchedule || form.getValues().cronSchedule == ''
      ? 'Not currently scheduled'
      : '';

  return (
    <Card>
      <CardHeader>
        <CardTitle>Schedule</CardTitle>
      </CardHeader>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <CardContent>
            <FormField
              control={form.control}
              name="cronSchedule"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Cron Schedule" {...field} />
                  </FormControl>
                  <FormDescription>
                    The schedule to run the job against if not a oneoff.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
          <CardFooter className="bg-muted">
            <div className="flex flex-row items-center justify-between w-full mt-4">
              <p className="text-muted-foreground text-sm">{msg}</p>
              <Button type="submit">Save</Button>
            </div>
          </CardFooter>
        </form>
      </Form>
    </Card>
  );
}

async function updateJobSchedule(
  jobId: string,
  cronSchedule?: string
): Promise<UpdateJobScheduleResponse> {
  const res = await fetch(`/api/jobs/${jobId}/schedule`, {
    method: 'PUT',
    headers: {
      'content-type': 'application/json',
    },
    body: JSON.stringify(
      new UpdateJobScheduleRequest({
        id: jobId,
        cronSchedule: cronSchedule,
      })
    ),
  });
  if (!res.ok) {
    const body = await res.json();
    throw new Error(body.message);
  }
  return UpdateJobScheduleResponse.fromJson(await res.json());
}

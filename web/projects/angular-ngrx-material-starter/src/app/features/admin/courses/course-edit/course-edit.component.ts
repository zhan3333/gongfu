import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import {
  AdminApiService,
  ISampleCoach
} from '../../../../api/admin/admin-api.service';
import { NotificationService } from '../../../../core/notifications/notification.service';
import { ActivatedRoute } from '@angular/router';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { School } from '../../../../api/models/school';
import * as moment from 'moment';
import { Course } from '../../../../api/models/course';
import { ApiService } from '../../../../api/api.service';
import { BottomSheetComponent } from '../../../../shared/bottom-sheet.component';
import { MatBottomSheet } from '@angular/material/bottom-sheet';

@Component({
  selector: 'anms-course-edit',
  templateUrl: './course-edit.component.html',
  styleUrls: ['./course-edit.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CourseEditComponent implements OnInit {
  public course: Course | undefined;
  public id = 0;

  // 教练数组
  public coaches: ISampleCoach[] = [];
  public loading = false;
  public schools: School[] = [];

  public form = new UntypedFormGroup({
    schoolId: new UntypedFormControl(null),
    startDate: new UntypedFormControl(null),
    startTime: new UntypedFormControl(null),
    managerId: new UntypedFormControl(null),
    coachId: new UntypedFormControl(null),
    assistantCoachIds: new UntypedFormControl([]),
    summary: new UntypedFormControl(''),
    content: new UntypedFormControl(''),
    images: new UntypedFormControl([]),
    checkInAt: new UntypedFormControl(''),
    checkOutAt: new UntypedFormControl(''),
    checkInBy: new UntypedFormControl(''),
    checkOutBy: new UntypedFormControl('')
  });
  private uploadProgressValue = 0;

  constructor(
    private adminApi: AdminApiService,
    public api: ApiService,
    private notification: NotificationService,
    public route: ActivatedRoute,
    private bottomSheet: MatBottomSheet
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      this.id = parseInt(params['id'], 10);
      if (this.id === 0) {
        this.notification.error('invalid request params');
        return;
      }
      this.adminApi.getCoaches().subscribe((data) => (this.coaches = data));
      this.adminApi.getSchools().subscribe((data) => (this.schools = data));
      this.refreshCourse(this.id);
    });
  }

  refreshCourse(id: number) {
    this.adminApi.getCourse(id).subscribe((course) => {
      this.course = course;
      console.log('course', this.course);
      this.form.patchValue({
        schoolId: this.course.school.id,
        startDate: moment(this.course.startDate, 'YYYY/MM/DD').toDate(),
        startTime: this.course.startTime,
        managerId: this.course.manager.id,
        coachId: this.course.coach?.id,
        assistantCoachIds: this.course.assistantCoachIds,
        summary: this.course.summary,
        content: this.course.content,
        images: this.course.images || [],
        checkInAt: this.course.checkInAt,
        checkOutAt: this.course.checkOutAt,
        checkInBy: this.course.checkInByUser?.id,
        checkOutBy: this.course.checkOutByUser?.id
      });
    });
  }

  onSave() {
    if (!this.course?.id) {
      return;
    }
    this.loading = true;
    const data = this.form.value;
    console.log(this.form.value);
    this.adminApi
      .updateCourse(this.course.id, {
        startDate: moment(data['startDate']).format('YYYY/MM/DD'),
        startTime: data['startTime'],
        schoolId: data['schoolId'],
        managerId: data['managerId'],
        coachId: data['coachId'],
        assistantCoachIds: data['assistantCoachIds'],
        summary: data['summary'],
        content: data['content'],
        images: data['images']
      })
      .subscribe(
        () => {
          this.notification.success('保存成功');
          this.refreshCourse(this.id);
        },
        (err) => {
          this.notification.error('保存失败');
          console.error('保存失败', err);
        },
        () => (this.loading = false)
      );
  }

  // 上传图片
  public onImageSelected(event: Event) {
    if (event.target == null) {
      return;
    }
    const input = event.target as HTMLInputElement;
    if (input.files === null || input.files.length === 0) {
      return;
    }
    const file = input.files[0];

    this.loading = true;
    this.api
      .uploadFile(
        file,
        'course',
        (value) => (this.uploadProgressValue = value),
        (fileKey) => {
          this.form.patchValue({
            images: [...this.form.value['images'], fileKey]
          });
        }
      )
      .subscribe(
        () => {},
        (error) => {
          this.notification.error(
            '上传头像失败，请稍后重试: ' + JSON.stringify(error)
          );
          this.loading = false;
          this.uploadProgressValue = 0;
        }
      );
  }

  // 删除图片
  removeImage(image: string) {
    this.bottomSheet
      .open(BottomSheetComponent, {
        data: new Map<string, string>([
          ['ok', '删除'],
          ['cancel', '取消']
        ])
      })
      .afterDismissed()
      .subscribe((res) => {
        if (res === 'ok') {
          this.form.patchValue({
            images: [
              ...this.form.value['images'].filter(
                (item: string) => item !== image
              )
            ]
          });
        }
      });
  }
}

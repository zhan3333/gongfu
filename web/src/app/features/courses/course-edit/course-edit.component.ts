import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { Course } from '../../../api/models/course';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { ApiService } from '../../../api/api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ActivatedRoute } from '@angular/router';
import { BottomSheetComponent } from '../../../shared/bottom-sheet.component';
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
  public uploadProgressValue = 0;
  public loading = false;
  public form = new UntypedFormGroup({
    managerId: new UntypedFormControl(null),
    coachId: new UntypedFormControl(null),
    assistantCoachIds: new UntypedFormControl([]),
    summary: new UntypedFormControl(''),
    content: new UntypedFormControl(''),
    images: new UntypedFormControl([])
  });

  constructor(
    public api: ApiService,
    private notification: NotificationService,
    private route: ActivatedRoute,
    private bottomSheet: MatBottomSheet
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      this.id = parseInt(params['id'], 10);
      if (this.id === 0) {
        this.notification.error('invalid request params');
        return;
      }
      this.refreshCourse(this.id);
    });
  }

  refreshCourse(id: number) {
    this.api.getCourse(id).subscribe((course) => {
      this.course = course;
      console.log('course', this.course);
      this.form.patchValue({
        managerId: this.course.manager.id,
        coachId: this.course.coach?.id,
        assistantCoachIds: this.course.assistantCoaches.map((item) => item.id),
        summary: this.course.summary,
        images: this.course.images || [],
        content: this.course.content
      });
    });
  }

  onSave() {
    console.log('coach', this.form.value);
    if (!this.course?.id) {
      return;
    }
    this.loading = true;
    const data = this.form.value;
    console.log(this.form.value);
    this.api
      .updateCourse(this.course.id, {
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

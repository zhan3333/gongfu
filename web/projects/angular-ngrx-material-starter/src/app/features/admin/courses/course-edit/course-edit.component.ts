import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { AdminApiService } from '../../../../api/admin/admin-api.service';
import { NotificationService } from '../../../../core/notifications/notification.service';
import { ActivatedRoute } from '@angular/router';
import { Course } from '../../../../api/models/Course';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'anms-course-edit',
  templateUrl: './course-edit.component.html',
  styleUrls: ['./course-edit.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CourseEditComponent implements OnInit {
  public course: Course | undefined;
  public form = new FormGroup({
    id: new FormControl(0),
    name: new FormControl(''),
    address: new FormControl(''),
    schoolStartAt: new FormControl(''),
    coachId: new FormControl(null),
    assistantCoachIds: new FormControl([]),
    summary: new FormControl(''),
    images: new FormControl([]),
    checkInAt: new FormControl(''),
    checkOutAt: new FormControl(''),
    checkInBy: new FormControl(''),
    checkOutBy: new FormControl(''),
  })

  constructor(
    private adminApi: AdminApiService,
    private notification: NotificationService,
    public route: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      const id = parseInt(params['id'], 10)
      if (id === 0) {
        this.notification.error('invalid request params')
        return
      }
      this.adminApi.getCourse(id).subscribe(course => {
        this.course = course
        console.log(this.course)
      })
    })
  }

}

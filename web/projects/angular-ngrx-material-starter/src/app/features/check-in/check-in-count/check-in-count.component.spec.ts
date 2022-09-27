import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckInCountComponent } from './check-in-count.component';

describe('CheckInCountComponent', () => {
  let component: CheckInCountComponent;
  let fixture: ComponentFixture<CheckInCountComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CheckInCountComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckInCountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

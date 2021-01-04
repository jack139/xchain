# coding:utf-8
import sys
import urllib3, json, base64, time, hashlib
from datetime import datetime

urllib3.disable_warnings()


# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys())
    if 'data' in name_list: # data 按 key 排序
        param['data'] = json.dumps(param['data'], sort_keys=True).replace(' ','')
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])


if __name__ == '__main__':
    if len(sys.argv)<2:
        print("usage: python3 %s <host> <port>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    port = sys.argv[2]

    body = {
        'version'  : '1',
        'sign_type' : 'SHA256', 
        'data'     : {
            'user_id'    : "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=",
            'data'     : 'zzzzz',
            'block_id' : '59534f7d-db5b-4792-8937-09996638c3d4',
            'action' : 1,
            'deal_id' : '59534f7d-db5b-4792-8937-09996638c3d4',
            'from_user_id' : 'j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=',
            'auth_id' : 'ef57fd9e-66c8-4d23-b142-8bc32b57bfcd',
        }
    }

    secret = 'UgM13IPx/BkwfQo8jceLq1CiXlT3lm4WLZ6K6TMR5bRueBGgTDVAv7ZLdBooTZWm2ixLaNitCW91NHW06h8VQw=='
    appid = hashlib.md5(secret.encode('utf-8')).hexdigest()
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appid'] = appid

    param_str = gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, secret)

    if body['sign_type'] == 'SHA256':
        sha256 = hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')
        signature_str =  base64.b64encode(sha256).decode('utf-8')
    else: # SM2
        #signature_str = sm2.SM2withSM3_sign_base64(sign_str)
        pass

    #print(sign_str)
    #print(sha256)
    #print(signature_str)

    body['sign_data'] = signature_str

    body = json.dumps(body)
    print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    host = 'http://%s:%s'%(hostname, port)
    url = host+'/api/query_block'

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
